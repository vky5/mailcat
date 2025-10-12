package imap

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/vky5/mailcat/internal/utils"

	"github.com/vky5/mailcat/internal/db/models"
)

// connect to the IMAP server
func ConnectIMAP(acc models.Account) (*client.Client, error) {
	// connection to IMAP server
	var conn *client.Client
	var err error

	address := fmt.Sprintf("%s:%s", acc.Host, acc.Port)

	// connect to the server and secure determines over tls or not
	// connect using TLS
	if acc.Secure {
		conn, err = client.DialTLS(address, &tls.Config{
			ServerName: acc.Host, // needed for TLS verification
		})
	} else {
		conn, err = client.Dial(address)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	// once the connection is established with acc.host then login
	// Login
	if err := conn.Login(acc.Email, acc.Password); err != nil {
		conn.Logout() // if login fails exit
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	log.Println("Connected and logged in to", acc.Host)

	return conn, err
}

// fetch emails
func FetchEmails(conn *client.Client, pageSize int, pageNumber int) ([]models.Email, error) {
	// select INBOX
	mbox, err := conn.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %v", err)
	}

	if mbox.Messages == 0 { // this means that there are 0 emails in INBOX
		log.Println("No messages found in INBOX")
		return nil, nil
	}

	from, to := utils.Paginate(int(mbox.Messages), pageSize, pageNumber) // we get the from and to values of the emails to be fetched

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(from), uint32(to)) // critical thing about this the order is guranteed.. meaning if 41th to 50 emails are fetched they are fetched in this sequence into the channel only

	// prepare the channel to fetch emails of pageSize (that is the number of emails to be fetched at a certain time)
	messages := make(chan *imap.Message, pageSize)
	done := make(chan error, 1)

	itmes := []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody} // the envelope contains from, to subject date etc and body contains attachment and actual body
	go func() {
		done <- conn.Fetch(seqset, itmes, messages) // read the emails and fill it in channel of buffer size pageSize where buffers are of type imap.Message that points to that in memory
		// once all emails have been fetched then it sends either nil or err to done and channel is closed by go-imap
	}()

	var emails []models.Email

	// how go knows when this channel is open or closed when reaading is done
	for msg := range messages {
		/*
			Go keeps reading from the channel until the channel is closed.
			Once the channel is closed and all buffered items have been read, the loop automatically exits.
			You don’t need to manually signal anything inside the loop.
		*/
		emails = parseMails(msg, emails)
	}
	// IMAP Server → conn.Fetch() → messages channel → for loop

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %v", err)
	}

	return emails, nil
}

// constantly listen to new emails
func ListenNewEmails() {

}

func Logout(conn *client.Client) {
	if err := conn.Logout(); err != nil {
		log.Println("Error logging out:", err)
	} else {
		log.Println("Logged out Successfully")
	}
}
