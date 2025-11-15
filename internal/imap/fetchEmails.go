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

	// Select mailbox
	mbox, err := conn.Select(mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select %s: %v", mailbox, err)
	}

	if mbox.Messages == 0 {
		return []models.Email{}, nil
	}

	// Pagination
	from, to := utils.Paginate(int(mbox.Messages), pageSize, pageNumber)

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(from), uint32(to))

	// Request:
	// - Envelope (meta)
	// - BodyStructure (parts)
	// - BODY[] full RFC822 message
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchBodyStructure,
		section.FetchItem(),
	}

	messages := make(chan *imap.Message, pageSize)
	done := make(chan error, 1)

	go func() {
		done <- conn.Fetch(seqset, items, messages)
	}()

	var emails []models.Email

	for msg := range messages {
		emails = parseMails(msg, emails)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %v", err)
	}

	utils.ReverseSlice(emails)
	return emails, nil
}
