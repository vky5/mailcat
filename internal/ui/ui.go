package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview" // high level package for TUI
	"github.com/vky5/mailcat/internal/db/models"
)

// StartUI builds the overall layout and starts the TUI.
func StartUI(accounts []*Account) error {
	app := tview.NewApplication()

	// ===== Middle Panel =====
	emailPanel := NewEmailListPanel(func(email models.Email) {
		// when user selects an email row
		// TODO: show email body in right panel
	})

	// ===== Right Panel =====
	placeholderRight := tview.NewBox().SetBorder(true).SetTitle("Content (placeholder)")

	// ===== Folder Selection Callback =====
	onSelect := func(account, folder string) {
		// Find selected folder and load emails
		for _, acc := range accounts {
			if acc.Email == account {
				for _, f := range acc.Folders {
					if f.Name == folder {
						emailPanel.SetEmails(f.Emails)
						break
					}
				}
			}
		}
	}

	// ===== Left Panel =====
	fp := NewFolderPanel(onSelect)

	// Populate accounts for testing
	for _, acc := range accounts {
		fp.AddAccount(acc.Email, acc.Folders)
	}

	// ===== Layout =====
	layout := tview.NewFlex().
		AddItem(fp.list, 30, 1, true).                // left has initial focus
		AddItem(emailPanel.Primitive(), 0, 2, false). // middle
		AddItem(placeholderRight, 0, 3, false)        // right

	// global keybindings to switch panels
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			if app.GetFocus() == fp.list {
				app.SetFocus(emailPanel.Primitive())
			} else if app.GetFocus() == emailPanel.Primitive() {
				app.SetFocus(placeholderRight)
			}
			return nil
		case tcell.KeyLeft:
			if app.GetFocus() == placeholderRight {
				app.SetFocus(emailPanel.Primitive())
			} else if app.GetFocus() == emailPanel.Primitive() {
				app.SetFocus(fp.list)
			}
			return nil
		}
		return event
	})

	return app.SetRoot(layout, true).Run()
}
