package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/db/models"
)

// StartUI builds the overall layout and starts the TUI.
func StartUI(accounts []*Account) error {
	app := tview.NewApplication()

	// ===== Right Panel =====
	emailOpenPanel := NewEmailOpenPanel()

	// ===== Middle Panel =====
	emailPanel := NewEmailListPanel(func(email models.Email) {
		// when user selects an email row, show it in the right panel
		emailOpenPanel.SetEmail(email)
		app.SetFocus(emailOpenPanel.Primitive())
	})

	// ===== Folder Selection Callback =====
	onSelect := func(account, folder string) {
		// Find selected folder and load emails
		for _, acc := range accounts {
			if acc.Email == account {
				for _, f := range acc.Folders {
					if f.Name == folder {
						emailPanel.SetEmails(f.Emails)
						// reset right panel to placeholder
						emailOpenPanel.Clear()
						// focus middle panel automatically
						app.SetFocus(emailPanel.Primitive())
						return
					}
				}
			}
		}
	}

	// ===== Left Panel =====
	fp := NewFolderPanel(onSelect)

	// Populate accounts
	for _, acc := range accounts {
		fp.AddAccount(acc.Email, acc.Folders)
	}

	// ===== Layout =====
	layout := tview.NewFlex().
		AddItem(fp.Primitive(), 30, 1, true).               // Left panel, initial focus
		AddItem(emailPanel.Primitive(), 0, 2, false).       // Middle panel
		AddItem(emailOpenPanel.Primitive(), 0, 3, false)    // Right panel

	// Global keybindings to switch panels
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			switch app.GetFocus() {
			case fp.Primitive():
				app.SetFocus(emailPanel.Primitive())
			case emailPanel.Primitive():
				app.SetFocus(emailOpenPanel.Primitive())
			}
			return nil
		case tcell.KeyLeft:
			switch app.GetFocus() {
			case emailOpenPanel.Primitive():
				app.SetFocus(emailPanel.Primitive())
			case emailPanel.Primitive():
				app.SetFocus(fp.Primitive())
			}
			return nil
		}
		return event
	})

	return app.SetRoot(layout, true).Run()
}
