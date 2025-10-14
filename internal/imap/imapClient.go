package imap

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/emersion/go-imap/client"
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

func Logout(conn *client.Client) {
	if err := conn.Logout(); err != nil {
		log.Println("Error logging out:", err)
	} else {
		log.Println("Logged out Successfully")
	}
}
