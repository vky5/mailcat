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
	list     *tview.List                  // actual visual component
	accounts []*Account                   // our data model
	onSelect func(account, folder string) // callback when a folder is selected
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

// getFolderIcon returns appropriate emoji for folder type
func getFolderIcon(folderName string) string {
	switch folderName {
	case "Inbox":
		return "📥"
	case "Sent":
		return "📤"
	case "Drafts":
		return "📝"
	case "Trash":
		return "🗑️"
	case "Spam":
		return "🚫"
	case "Archive":
		return "📦"
	case "Important":
		return "⭐"
	default:
		return "📁"
	}
}

// describing the basic layout of the panel returning the folder panel
func NewFolderPanel(onSelect func(account string, folder string)) *FolderPanel {
	fp := &FolderPanel{
		list:     tview.NewList(),
		onSelect: onSelect,
	}

	fp.list.SetBorder(true).SetTitle(" 📂 Accounts ")
	fp.list.SetBorderColor(tcell.NewRGBColor(0, 191, 255))    // bright cyan to match email panel
	fp.list.SetBackgroundColor(tcell.NewRGBColor(18, 30, 40)) // dark background
	fp.list.ShowSecondaryText(false)
	fp.list.SetHighlightFullLine(true)
	fp.list.SetMainTextColor(tcell.NewRGBColor(180, 220, 255))         // soft cyan for text
	fp.list.SetSelectedBackgroundColor(tcell.NewRGBColor(0, 100, 150)) // match email panel selection
	fp.list.SetSelectedTextColor(tcell.ColorWhite)

	return fp
}

// helper func to add account to the panel
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

	// Add New Account with vibrant styling
	fp.list.AddItem("[::b][#32CD32]✨ Add New Account[-:-:-]", "", 0, func() {
		fmt.Println("Add account clicked")
	})

	// Add separator (non-selectable)
	fp.list.AddItem("[#2F4F4F]────────────────────────[-]", "", 0, nil)

	for accIdx, acc := range fp.accounts {
		// Top-level account with email icon and better styling
		totalUnread := 0
		for _, f := range acc.Folders {
			totalUnread += f.UnreadCount()
		}

		expandIcon := "▶"
		if acc.Expanded {
			expandIcon = "▼"
		}

		accText := fmt.Sprintf("[::b][#00BFFF]%s 📧 %s[-:-:-]", expandIcon, acc.Email)
		if totalUnread > 0 {
			accText = fmt.Sprintf("[::b][#00BFFF]%s 📧 %s [#FFD700](%d)[-:-:-]", expandIcon, acc.Email, totalUnread)
		}

		fp.list.AddItem(accText, "", 0, func(a *Account) func() {
			return func() {
				a.Expanded = !a.Expanded
				fp.render()
			}
		}(acc))

		if acc.Expanded {
			// Folder children with better icons and styling
			for i, f := range acc.Folders {
				line := "├──"
				connector := "│"
				if i == len(acc.Folders)-1 {
					line = "└──"
					connector = " "
				}

				unread := f.UnreadCount()
				folderIcon := getFolderIcon(f.Name)
				
				// Base folder color
				folderColor := "#B0B0B0" // gray for read
				folderStyle := ""
				
				// Highlight if unread
				if unread > 0 {
					folderColor = "#FFD700" // gold for unread
					folderStyle = "::b"
				}

				// Format unread count
				unreadText := ""
				if unread > 0 {
					unreadText = fmt.Sprintf(" [#32CD32::b](%d)[-:-:-]", unread)
				} else {
					unreadText = fmt.Sprintf(" [#778899](%d)[-]", unread)
				}

				folderText := fmt.Sprintf("  [#4682B4]%s[-] [%s%s]%s %s[-:-:-]%s", 
					line, 
					folderColor, 
					folderStyle,
					folderIcon, 
					f.Name,
					unreadText,
				)

				folder := f // copy for closure
				accountEmail := acc.Email
				fp.list.AddItem(folderText, "", 0, func() {
					fp.onSelect(accountEmail, folder.Name)
				})

				// Add subtle connector line between folders (non-selectable)
				if i < len(acc.Folders)-1 {
					connectorText := fmt.Sprintf("  [#4682B4]%s[-]", connector)
					fp.list.AddItem(connectorText, "", 0, nil)
				}
			}
			
			// Add separator after expanded account (non-selectable)
			if accIdx < len(fp.accounts)-1 {
				fp.list.AddItem("[#2F4F4F]────────────────────────[-]", "", 0, nil)
			}
		}
	}
}

// Primitive returns the tview primitive
func (fp *FolderPanel) Primitive() tview.Primitive {
	return fp.list
}
