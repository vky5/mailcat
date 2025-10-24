package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/db/models"
)

type Folder struct {
	Name   string // e.g Inbox, sent, trash
	Emails []models.Email
}

// Account is emails
type Account struct {
	Email    string
	Folders  []Folder
	Expanded bool // shows whether we show the folders or not
}

// struct that connects data to UI
type FolderPanel struct {
	list     *tview.List                  // actual visual component // tview.list A list displays rows of text. The user can navigate them with arrow keys and select them
	accounts []*Account                   // our data model
	onSelect func(account, folder string) // callback when a folder is selected
	// when multiple parameters in a func have same type we can omit repeating type until it changes
}

// get unread emails count
func (f *Folder) UnreadCount() int {
	count := 0
	for _, e := range f.Emails {
		if !e.Read {
			count++
		}
	}
	return count
}

// describing the basic layout of the panel returning the folder pannel
func NewFolderPanel(onSelect func(account string, folder string)) *FolderPanel {
	fp := &FolderPanel{
		list:     tview.NewList(),
		onSelect: onSelect,
	}

	fp.list.SetBorder(true).SetTitle("Accounts")
	fp.list.SetBorderColor(tcell.ColorDarkGray)               // border color
	fp.list.SetBackgroundColor(tcell.NewRGBColor(18, 30, 40)) // dark background
	fp.list.ShowSecondaryText(false)
	fp.list.SetHighlightFullLine(true)
	fp.list.SetMainTextColor(tcell.NewRGBColor(180, 220, 255))       // soft cyan for text
	fp.list.SetSelectedBackgroundColor(tcell.NewRGBColor(0, 50, 70)) // muted blue for selection
	fp.list.SetSelectedTextColor(tcell.ColorWhite)

	// Add New Account with bright green
	fp.list.AddItem("[::b][#32CD32]+ Add New Account[-:-:-]", "", '+', nil)

	return fp
}

// helper func to add account to the pannel
func (fp *FolderPanel) AddAccount(email string, folders []Folder) {
	account := &Account{
		Email:    email,
		Folders:  folders,
		Expanded: false,
	}

	fp.accounts = append(fp.accounts, account)
	fp.render()
}

// refresh the panel with new emails or updates
func (fp *FolderPanel) render() {
	fp.list.Clear()

	// Add New Account (bright green, no background tag)
	fp.list.AddItem("[::b][#32CD32]+ Add New Account[-:-:-]", "", '+', func() {
		fmt.Println("Add account clicked")
	})

	for _, acc := range fp.accounts {
		// Top-level account (cyan, no background tag)
		accText := fmt.Sprintf("[::b][#00BFFF]%s[-:-:-]", acc.Email)

		fp.list.AddItem(accText, "", 0, func(a *Account) func() { // an IIFE which takes account as parameter to call another func immediately
			// capture pointer safely
			return func() {
				a.Expanded = !a.Expanded
				fp.render()
			}
		}(acc))

		if acc.Expanded {
			// Folder children with tree lines
			for i, f := range acc.Folders {
				line := "├──"
				if i == len(acc.Folders)-1 {
					line = "└──"
				}
				unread := f.UnreadCount()
				folderColor := "#B0B0B0" // gray
				if unread > 0 {
					folderColor = "#FFD700" // yellow for unread
				}
				folderText := fmt.Sprintf("  [%s]%s %s (%d)[-]", folderColor, line, f.Name, unread)

				folder := f // copy for closure
				accountEmail := acc.Email
				fp.list.AddItem(folderText, "", 0, func() {
					fp.onSelect(accountEmail, folder.Name)
				})
			}
			fp.list.AddItem("", "", 0, nil)
		}
	}
}
