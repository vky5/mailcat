package main

import (
	"time"

	"github.com/vky5/mailcat/internal/db/models"
	"github.com/vky5/mailcat/internal/logger"
	"github.com/vky5/mailcat/internal/ui"
)

func main() {

	// setting up logger
	err := logger.Init("mailcat", false)
	if err != nil {
		panic(err)
	}

	// ===== Fake Data Setup =====
	inboxEmails := []models.Email{
		{
			ID:          1,
			From:        "alice@example.com",
			Subject:     "Meeting Reminder",
			Body:        "Don't forget about our meeting at 3PM.",
			Date:        time.Date(2025, 10, 24, 15, 0, 0, 0, time.Local),
			Read:        false,
			Attachments: "fafsdfsdf",
		},
		{
			ID:      2,
			From:    "newsletter@cloudflare.com",
			Subject: "Cloudflare Workers Weekly Update",
			Body:    "New features, community tools, and upcoming talks.",
			Date:    time.Date(2025, 10, 23, 10, 0, 0, 0, time.Local),
			Read:    true,
		},
		{
			ID:      3,
			From:    "bob@example.com",
			Subject: "Lunch Tomorrow?",
			Body:    "Hey, want to grab lunch tomorrow?",
			Date:    time.Date(2025, 10, 22, 12, 30, 0, 0, time.Local),
			Read:    false,
		},
	}

	sentEmails := []models.Email{
		{
			ID:      4,
			From:    "you@example.com",
			Subject: "Follow-up: Project Proposal",
			Body:    "Hey, just following up on the project proposal...",
			Date:    time.Date(2025, 10, 20, 17, 45, 0, 0, time.Local),
			Read:    true,
		},
	}

	trashEmails := []models.Email{
		{
			ID:      5,
			From:    "spam@random.com",
			Subject: "You won a prize!",
			Body:    "Click here to claim your reward.",
			Date:    time.Date(2025, 10, 19, 9, 30, 0, 0, time.Local),
			Read:    true,
		},
	}

	// ===== Accounts =====
	accounts := []*ui.Account{
		{
			Email: "vaibhavk@a.com",
			Folders: []ui.Folder{
				{Name: "Inbox", Emails: inboxEmails},
				{Name: "Sent", Emails: sentEmails},
				{Name: "Trash", Emails: trashEmails},
			},
		},
		{
			Email: "test@domain.com",
			Folders: []ui.Folder{
				{Name: "Inbox", Emails: inboxEmails},
				{Name: "Archive", Emails: sentEmails},
			},
		},
	}

	// ===== Launch UI =====
	if err := ui.StartUI(accounts); err != nil {
		panic(err)
	}
}
