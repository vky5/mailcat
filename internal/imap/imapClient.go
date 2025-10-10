package imap

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/vky5/mailcat/internal/utils"
)

// connection to IMAP server
var conn *client.Client
var err error

// Accounts credentials for IMAP connection
type Account struct {
	Host   string
	Port   string
	Secure bool
	User   string
	Pass   string
}

// connect to the IMAP server
func ConnectIMAP(acc Account) (*client.Client, error) {

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
	if err := conn.Login(acc.User, acc.Pass); err != nil {
		conn.Logout() // if login fails exit
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	log.Println("Connected and logged in to", acc.Host)

	return conn, err
}

// fetch emails
func FetchEmails(pageSize int, pageNumber int) ([]Email, error) {
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
	seqset.AddRange(uint32(from), uint32(to))

	// prepare the channel to fetch emails of pageSize (that is the number of emails to be fetched at a certain time)
	messages := make(chan *imap.Message, pageSize)
	done := make(chan error, 1)

	itmes := []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody} // the envelope contains from, to subject date etc and body contains attachment and actual body 
	go func() {
		done <- conn.Fetch(seqset, itmes, messages)
	}()

	var emails []Email

	for msg := range messages {
		emails = parseMails(msg, emails)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %v", err)
	}

	return emails, nil
}

// constantly listen to new emails
func ListenNewEmails() {

}

func Logout() {
	if err := conn.Logout(); err != nil {
		log.Println("Error logging out:", err)
	} else {
		log.Println("Logged out Successfully")
	}
}
