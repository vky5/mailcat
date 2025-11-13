package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/commands"
	"github.com/vky5/mailcat/internal/db/models"
	"github.com/vky5/mailcat/internal/logger"
)

// StartUI builds the overall layout and starts the TUI.
func StartUI(accounts []*Account) error {
	app := tview.NewApplication()

	// ===== Right Panel =====
	emailOpenPanel := NewEmailOpenPanel()

	// ===== Middle Panel =====
	emailPanel := NewEmailListPanel(func(email models.Email) {
		emailOpenPanel.SetEmail(email)
		app.SetFocus(emailOpenPanel.Primitive())
	})

	// ===== Folder Selection Callback =====
	onSelect := func(account, folder string) {
		for _, acc := range accounts {
			if acc.Email == account {
				for _, f := range acc.Folders {
					if f.Name == folder {
						emailPanel.SetEmails(f.Emails)
						emailOpenPanel.Clear()
						app.SetFocus(emailPanel.Primitive())
						return
					}
				}
			}
		}
	}

	// ===== Left Panel =====
	fp := NewFolderPanel(onSelect)
	for _, acc := range accounts {
		fp.AddAccount(acc.Email, acc.Folders)
	}

	// ===== Command Bar =====
	cmdBar := NewCommandBar(app)

	// registering commands
	helpCmd := commands.NewHelpCommand(cmdBar.registry)
	cmdBar.Register(helpCmd)
	cmdBar.Register(&commands.AddAccount{})

	// ===== Layout =====
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	upperLayout := tview.NewFlex().
		AddItem(fp.Primitive(), 30, 1, true).
		AddItem(emailPanel.Primitive(), 0, 2, false).
		AddItem(emailOpenPanel.Primitive(), 0, 3, false)

	mainLayout.AddItem(upperLayout, 0, 1, true)
	mainLayout.AddItem(cmdBar.GetPrimitive(), 3, 0, false)

	// Track last focused panel
	var lastFocus tview.Primitive = fp.Primitive()

	// Global keybindings
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			switch app.GetFocus() {
			case fp.Primitive():
				app.SetFocus(emailPanel.Primitive())
				lastFocus = emailPanel.Primitive()
			case emailPanel.Primitive():
				app.SetFocus(emailOpenPanel.Primitive())
				lastFocus = emailOpenPanel.Primitive()
			}
			return nil
		case tcell.KeyLeft:
			switch app.GetFocus() {
			case emailOpenPanel.Primitive():
				app.SetFocus(emailPanel.Primitive())
				lastFocus = emailPanel.Primitive()
			case emailPanel.Primitive():
				app.SetFocus(fp.Primitive())
				lastFocus = fp.Primitive()
			}
			return nil
		}

		// Focus command bar
		if event.Key() == tcell.KeyRune && event.Rune() == ':' {
			lastFocus = app.GetFocus()
			cmdBar.input.SetText("")
			app.SetFocus(cmdBar.input)
			return nil
		}

		return event
	})

	// Command bar done function
	cmdBar.input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			input := cmdBar.input.GetText()
			logger.Info("Command entered: ", input)
			cmdBar.handleInput(input)

			// If no active command is running, restore focus
			if cmdBar.active == nil {
				app.SetFocus(lastFocus)
			}

		case tcell.KeyEsc:
			cmdBar.input.SetText("")
			cmdBar.ShowMessage("Type !help for commands")
			app.SetFocus(lastFocus)
		}
	})

	return app.SetRoot(mainLayout, true).Run()
}
