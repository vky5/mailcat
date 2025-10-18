package ui

import (
	"github.com/gdamore/tcell/v2"
	// "github.com/prometheus/alertmanager/notify/email"
	"github.com/rivo/tview"
)

type Folder struct {
	Name   string // e.g Inbox, sent, trash
	Unread int
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
	// when multiple parameters in a func have same type we can omit repeating type until it changes
}

func NewFolderPanel(onSelect func(account string, folder string)) *FolderPanel {
	fp := &FolderPanel{
		list:     tview.NewList(),
		onSelect: onSelect,
	}

	fp.list.SetBorder(true).SetTitle("Accounts")
	fp.list.SetBackgroundColor(tcell.NewRGBColor(0, 70, 70))
	fp.list.ShowSecondaryText(false)
	fp.list.SetHighlightFullLine(true)
	fp.list.AddItem("+ Add New Account", "", '+', nil)

	return fp
}

func (fp *FolderPanel) AddAccount(email string, folders []Folder) {
	account := &Account{
		Email:    email,
		Folders:  folders,
		Expanded: false,
	}

	fp.accounts = append(fp.accounts, account)
	fp.render()

}

func (fp *FolderPanel) render() {
	fp.list.Clear()
	fp.list.AddItem("+ Add New Account", "", '+', nil)

	for _, account := range fp.accounts {
		fp.list.AddItem(account.Email, "", 0, nil)
	}
}
