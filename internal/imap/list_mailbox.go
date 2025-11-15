package imap

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)


func ListMailboxes(conn *client.Client) ([]string, error) {
	mailboxes := []string{}

	// "" = root, "*" = wildcard (all mailboxes)
	mboxChan := make(chan *imap.MailboxInfo, 20)
	done := make(chan error, 1)

	go func() {
		done <- conn.List("", "*", mboxChan)
	}()

	for m := range mboxChan {
		mailboxes = append(mailboxes, m.Name)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("error listing mailboxes: %v", err)
	}

	return mailboxes, nil
}