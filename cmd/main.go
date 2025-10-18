package main

import (
	"github.com/vky5/mailcat/internal/ui"
)

func main() {
	// Mock accounts and folders
	accounts := []*ui.Account{
		{
			Email: "vaibhavk@a.com",
			Folders: []ui.Folder{
				{Name: "Inbox"},
				{Name: "Sent"},
				{Name: "Trash"},
			},
		},
		{
			Email: "test@domain.com",
			Folders: []ui.Folder{
				{Name: "Inbox"},
				{Name: "Archive"},
			},
		},
	}

	// Start terminal UI
	if err := ui.StartUI(accounts); err != nil {
		panic(err)
	}
}
