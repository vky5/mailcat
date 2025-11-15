package main

import (

	"github.com/vky5/mailcat/internal/db"
	"github.com/vky5/mailcat/internal/db/models"
	"github.com/vky5/mailcat/internal/imap"
	"github.com/vky5/mailcat/internal/logger"
	"github.com/vky5/mailcat/internal/ui"
)

func main() {
	// setting up logger
	err := logger.Init("mailcat.log", false)
	if err != nil {
		panic(err)
	}

	db.InitDB()

	var dbAccounts []models.Account
	if err := db.DB.Find(&dbAccounts).Error; err != nil {
		logger.Log.Fatalf("Failed to fetch accounts from DB: %v", err)
	}

	// convert DB models to UI account pointers expected by StartUI
	uiAccounts := make([]*ui.Account, len(dbAccounts))

	for i := range dbAccounts {
		conn, err := imap.GetConnection(dbAccounts[i])
		if err != nil {
			logger.Error("Failed IMAP connection for", dbAccounts[i].Email, err)
			continue
		}

		folders, err := imap.ListMailboxes(conn)
		if err != nil {
			logger.Error("Failed to list mailboxes", err)
			folders = []string{"INBOX"}
		}

		uiFolders := make([]ui.Folder, len(folders))
		for j, f := range folders {
			uiFolders[j] = ui.Folder{Name: f}
			logger.Log.Println("Loaded folder:", f)
		}

		uiAccounts[i] = &ui.Account{
			ID:      dbAccounts[i].ID,
			Email:   dbAccounts[i].Email,
			Folders: uiFolders,
		}
	}

	// ===== Launch UI =====
	if err := ui.StartUI(uiAccounts); err != nil {
		panic(err)
	}
}
