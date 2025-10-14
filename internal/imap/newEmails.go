package imap

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap/client"
	"github.com/vky5/mailcat/internal/db/models"
)

// constantly listen to new emails on a given IMAP connnection and mailbox
func ListenNewEmails(conn *client.Client, mailbox string, out chan models.Email) error {
	mbox, err := conn.Select(mailbox, false)
	if err != nil {
		return fmt.Errorf("failed to select mailbox %s: %v", mailbox, err)
	}
	log.Printf("Listening for new emails in %s (currently %d messages)\n", mailbox, mbox.Messages)

	updates := make(chan client.Update) // type from go-imap that represents any kind of updates the IMAP server sends (new message, message deletion, flag change)
	conn.Updates = updates              // Updates is a conn's field which tells conn whenever the server sends any update, push it into this channel

	for {
		stop := make(chan struct{})
		done := make(chan error, 1)

		// sit idle and receive push style updates liek new message so instead of polling every few sec the client sends IDLE and server responds with idling
		// for new event server sends `23 EXISTS` to stop idling the client sends DONE
		go func() {
			done <- conn.Idle(stop, nil) // idle is to keeping the loop for lisening to new message open and if there is any new message or something else, it sends that update to updaate channel, the done is our own channel listening to current state of the client (us) like if it exited properly or something happened

			// first argument of Idle (stop) is the control argument to tell goroutine when to stop
			// the second argument can be the update channel if we wanted it to receive messages but we already have made this conn.Updates channel for delivering message
		}()

		select {
		case update := <-updates:
			// check for new messages
			if mboxUpdate, ok := update.(*client.MailboxUpdate); ok {
				if mboxUpdate.Mailbox.Messages > mbox.Messages {
					newCount := mboxUpdate.Mailbox.Messages - mbox.Messages
					log.Printf("New %d messages detected\n", newCount)
				}

				from := mbox.Messages + 1
				to := mboxUpdate.Mailbox.Messages
				emails, err := FetchEmails(conn, int(to-from+1), 1)
				if err != nil {
					log.Printf("Error fetching new emails: %v", err)
				} else {
					for _, e := range emails {
						out <- e
					}
				}

				mbox.Messages = mboxUpdate.Mailbox.Messages
			}

		// close after every 30 mins and restart
		case <-time.After(30 * time.Minute): // according to RFC standard, the client (our server) cant open the connection and just dissapear meaning it closes automatically after 30 mins and this is where this comes.
			// we close the connection using close(stop) and
			close(stop)
			<-done // drain done for nil or any err (but err not possible)

		}

		// stop IDLE (this is for safety purpose in case the select blocks encounters any error)
		close(stop) // closing channel stop
		if err := <-done; err != nil {
			log.Println("Idle error: ", err)
		}

	}
}
