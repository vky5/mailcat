package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/commands"
	"github.com/vky5/mailcat/internal/db"
	"github.com/vky5/mailcat/internal/db/models"
	"github.com/vky5/mailcat/internal/imap"
	"github.com/vky5/mailcat/internal/logger"
	"strings"
)

// StartUI builds the overall layout and starts the TUI.
func StartUI(_ []*Account) error {
	logger.Info("Starting UI initialization...")
	app := tview.NewApplication()

	// ===== Right Panel =====
	logger.Info("Creating email open panel...")
	emailOpenPanel := NewEmailOpenPanel()

	// ===== Middle Panel =====
	logger.Info("Creating email list panel...")
	emailPanel := NewEmailListPanel(func(email models.Email) {
		logger.Info("Email selected:", email.Subject)
		emailOpenPanel.SetEmail(email)
		app.SetFocus(emailOpenPanel.Primitive())
	})

	// loader to show when fetching emails
	// loader := NewLoader("Loading emails")

	// accounts kept in UI memory
	accounts := []*Account{}

	// ===== Folder Selection Callback =====
	onSelect := func(accountEmail, folderName string) {
		logger.Info("Folder selected - Account:", accountEmail, "Folder:", folderName)

		// find account
		var acc *Account
		for _, a := range accounts {
			if a.Email == accountEmail {
				acc = a
				break
			}
		}
		if acc == nil {
			logger.Error("Account not found:", accountEmail)
			return
		}

		logger.Info("Account found in UI memory, showing loader...")

		// show loader immediately (we're already in UI thread context)
		logger.Info("Setting loader...")
		emailPanel.SetLoading(NewLoader("Fetching emails..."))

		logger.Info("Starting goroutine for email fetch...")
		go func() {
			logger.Info("Goroutine started for folder:", folderName)

			// Helper to update UI with error
			showError := func(msg string, err error) {
				logger.Error(msg, err)
				app.QueueUpdateDraw(func() {
					logger.Info("QueueUpdateDraw: Clearing loading state due to error")
					emailPanel.SetEmails([]models.Email{}) // Clear loading state
					// Optionally show error message in the panel
				})
			}

			// load full DB model
			logger.Info("Fetching account from database:", acc.Email)
			var dbAcc models.Account
			if err := db.DB.Where("email = ?", acc.Email).First(&dbAcc).Error; err != nil {
				showError("DB fetch failed:", err)
				return
			}
			logger.Info("Database account fetched successfully")

			// IMAP connection with timeout context
			logger.Info("Attempting IMAP connection...")
			conn, err := imap.GetConnection(dbAcc)
			if err != nil {
				showError("IMAP reconnect failed:", err)
				return
			}
			logger.Info("IMAP connection established successfully")

			// clean mailbox name
			clean := strings.TrimSpace(folderName)
			clean = strings.Trim(clean, `"`)
			logger.Info("Cleaned folder name:", clean)

			logger.Info("Starting FetchEmails from:", clean)
			emails, err := imap.FetchEmails(conn, clean, 50, 1)
			if err != nil {
				showError("Failed fetching from "+folderName+":", err)
				return
			}

			logger.Info("Fetched", len(emails), "emails from", clean)

			// now update UI from UI-safe context
			logger.Info("Queueing UI update with fetched emails...")
			app.QueueUpdateDraw(func() {
				logger.Info("QueueUpdateDraw: Setting", len(emails), "emails")
				emailPanel.SetEmails(emails)
				emailOpenPanel.Clear()
				app.SetFocus(emailPanel.Primitive())
				logger.Info("UI updated successfully with emails")
			})
		}()
		logger.Info("Goroutine spawned, returning from onSelect")
	}

	// ===== Command Bar (must be created BEFORE folder panel) =====
	logger.Info("Creating command bar...")
	cmdBar := NewCommandBar(app)

	// register commands
	helpCmd := commands.NewHelpCommand(cmdBar.registry)
	cmdBar.Register(helpCmd)
	cmdBar.Register(commands.NewAddAccount())

	// ===== Left Panel (Folder List) =====
	logger.Info("Creating folder panel...")
	fp := NewFolderPanel(onSelect, func() {
		// Trigger !addaccount as if typed
		addCmd := cmdBar.registry["!addaccount"]
		cmdBar.active = addCmd
		cmdBar.ShowMessage("Running: " + addCmd.Description())
		addCmd.Begin(cmdBar)

		// Switch focus to command bar input
		app.SetFocus(cmdBar.input)
	})

	// load DB accounts
	logger.Info("Loading accounts from database...")
	var accountsList []models.Account
	if err := db.DB.Find(&accountsList).Error; err != nil {
		logger.Error("Failed to load accounts from database:", err)
		return err
	}
	logger.Info("Loaded", len(accountsList), "accounts from database")

	// Build left panel with REAL IMAP folders
	for _, acc := range accountsList {
		logger.Info("Processing account:", acc.Email)
		uiAcc := &Account{
			Email:   acc.Email,
			Folders: []Folder{{Name: "INBOX", Emails: []models.Email{}}}, // Default
		}

		logger.Info("Adding account to folder panel:", acc.Email)
		fp.AddAccount(acc.Email, uiAcc.Folders)
		accounts = append(accounts, uiAcc)

		// Load folders asynchronously
		logger.Info("Starting async folder load for:", acc.Email)
		go func(account models.Account, uiAccount *Account) {
			logger.Info("Async goroutine started for:", account.Email)

			logger.Info("Getting IMAP connection for folder list:", account.Email)
			conn, err := imap.GetConnection(account)
			if err != nil {
				logger.Error("IMAP login failed for", account.Email, ":", err)
				return
			}
			logger.Info("IMAP connection successful for:", account.Email)

			logger.Info("Listing mailboxes for:", account.Email)
			boxes, err := imap.ListMailboxes(conn)
			if err != nil {
				logger.Error("Folder fetch failed for", account.Email, ":", err)
				return
			}
			logger.Info("Found", len(boxes), "mailboxes for", account.Email)

			folders := make([]Folder, len(boxes))
			for i, name := range boxes {
				folders[i] = Folder{Name: name, Emails: []models.Email{}}
			}

			logger.Info("Queueing UI update to add folders for:", account.Email)
			app.QueueUpdateDraw(func() {
				logger.Info("QueueUpdateDraw: Updating folders for:", account.Email)
				uiAccount.Folders = folders
				fp.UpdateAccount(account.Email, folders)
				logger.Info("Folders updated successfully for:", account.Email)
			})
		}(acc, uiAcc)
	}

	// ===== Command Bar =====
	cmdBar.Register(helpCmd)
	cmdBar.Register(commands.NewAddAccount())

	// ===== Layout =====
	logger.Info("Building layout...")
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	upperLayout := tview.NewFlex().
		AddItem(fp.Primitive(), 30, 1, true).
		AddItem(emailPanel.Primitive(), 0, 2, false).
		AddItem(emailOpenPanel.Primitive(), 0, 3, false)

	mainLayout.AddItem(upperLayout, 0, 1, true)
	mainLayout.AddItem(cmdBar.GetPrimitive(), 3, 0, false)

	// navigation
	var lastFocus tview.Primitive = fp.Primitive()

	logger.Info("Setting up input capture...")
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			switch app.GetFocus() {
			case fp.Primitive():
				logger.Info("Navigation: Moving focus from folder panel to email list")
				app.SetFocus(emailPanel.Primitive())
				lastFocus = emailPanel.Primitive()
			case emailPanel.Primitive():
				logger.Info("Navigation: Moving focus from email list to email open panel")
				app.SetFocus(emailOpenPanel.Primitive())
				lastFocus = emailOpenPanel.Primitive()
			}
			return nil
		case tcell.KeyLeft:
			switch app.GetFocus() {
			case emailOpenPanel.Primitive():
				logger.Info("Navigation: Moving focus from email open panel to email list")
				app.SetFocus(emailPanel.Primitive())
				lastFocus = emailPanel.Primitive()
			case emailPanel.Primitive():
				logger.Info("Navigation: Moving focus from email list to folder panel")
				app.SetFocus(fp.Primitive())
				lastFocus = fp.Primitive()
			}
			return nil
		}

		if event.Key() == tcell.KeyRune && event.Rune() == ':' {
			logger.Info("Command mode activated")
			lastFocus = app.GetFocus()
			cmdBar.input.SetText("")
			app.SetFocus(cmdBar.input)
			return nil
		}

		return event
	})

	cmdBar.input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			input := cmdBar.input.GetText()
			logger.Info("Command entered:", input)
			cmdBar.handleInput(input)

			if cmdBar.active == nil {
				logger.Info("Returning focus to last focused panel")
				app.SetFocus(lastFocus)
			}

		case tcell.KeyEsc:
			logger.Info("Command mode cancelled")
			cmdBar.input.SetText("")
			cmdBar.ShowMessage("Type !help for commands")
			app.SetFocus(lastFocus)
		}
	})

	logger.Info("Starting TUI application...")
	return app.SetRoot(mainLayout, true).Run()
}
