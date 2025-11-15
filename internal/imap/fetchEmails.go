package imap

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/vky5/mailcat/internal/db/models"
	"github.com/vky5/mailcat/internal/utils"
)

// FetchEmails returns a paginated list of emails from a mailbox.
func FetchEmails(conn *client.Client, mailbox string, pageSize int, pageNumber int) ([]models.Email, error) {

	// Select the mailbox. Gmail does NOT allow selecting folders like "[Gmail]".
	mbox, err := conn.Select(mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select %s: %v", mailbox, err)
	}

	// No messages? Return empty slice.
	if mbox.Messages == 0 {
		return []models.Email{}, nil
	}

	// Compute "from" and "to" sequence numbers for pagination.
	from, to := utils.Paginate(int(mbox.Messages), pageSize, pageNumber)

	// Prepare sequence set for IMAP (e.g. 41–50)
	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(from), uint32(to))

	// IMPORTANT: Gmail rejects imap.FetchBody (no section).
	// We must explicitly request BODY[] section.
	section := &imap.BodySectionName{}

	items := []imap.FetchItem{
		imap.FetchEnvelope,      // From, To, Subject, Date
		section.FetchItem(),     // BODY[] = safe body fetch compatible with Gmail
	}

	// IMAP messages are streamed into this channel.
	messages := make(chan *imap.Message, pageSize)
	done := make(chan error, 1)

	// Run the fetch in a goroutine.
	go func() {
		done <- conn.Fetch(seqset, items, messages)
	}()

	var emails []models.Email

	// Read all messages until channel closes.
	for msg := range messages {
		// parseMails fills your models.Email struct
		emails = parseMails(msg, emails)
	}

	// Check if Fetch returned an error.
	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %v", err)
	}

	// IMAP returns oldest → newest, so reverse for UI.
	utils.ReverseSlice(emails)

	return emails, nil
}
