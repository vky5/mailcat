package ui

import (
	// "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func StartUI(accounts []*Account) error {
	app := tview.NewApplication()

	// Callback when folder is selected
	onSelect := func(account, folder string) {
		// For now, just print selection to console
		// fmt.Printf("Selected: %s -> %s\n", account, folder)
	}

	// ===== Left Panel =====
	fp := NewFolderPanel(onSelect)

	// Populate accounts for testing
	for _, acc := range accounts {
		fp.AddAccount(acc.Email, acc.Folders)
	}

	// ===== Placeholder Panels =====
	placeholderMiddle := tview.NewBox().SetBorder(true).SetTitle("Emails (placeholder)")
	placeholderRight := tview.NewBox().SetBorder(true).SetTitle("Content (placeholder)")

	// ===== Layout =====
	layout := tview.NewFlex().
		AddItem(fp.list, 30, 1, true).
		AddItem(placeholderMiddle, 0, 2, false).
		AddItem(placeholderRight, 0, 3, false)

	return app.SetRoot(layout, true).Run()
}
