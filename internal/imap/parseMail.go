package imap

import "time"

type Email struct {
	From        string
	To          []string
	Subject     string
	Body        string
	Date        time.Time	
}


