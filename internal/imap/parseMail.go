package imap

import (
	"io"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"github.com/vky5/mailcat/internal/db/models"
)

func parseMails(msg *imap.Message, emails []models.Email) []models.Email { // using pointer to save memory
	section := &imap.BodySectionName{} // give the entire body fo the message including headers and texts
	// section := &imap.BodySectionName{Peek: true, Path: []string{"TEXT"}} // if u want only text part // msg.GetBody(section)

	r := msg.GetBody(section)
	if r == nil {
		return emails
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return emails
	}

	var body string

	// loop through all message parts
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		// each part from mr.NextPart() has two things
		// part.Header -> meta data (like content type content disposition etc)
		// part.Body -> actual data (like email text or an attachment stream)

		// part.Header.(type) here type depends on what the actual data is liek mail.InlineHeader mail.AttachemntHeader
		switch part.Header.(type) {
		case *mail.InlineHeader: // if MIME part is inline content (like main body text) read its body stream completly into memory and convert it to string
			b, _ := io.ReadAll(part.Body)
			body = string(b)
		}
	}

	// Extract header info
	var from string
	var to []string
	var subject string
	var date time.Time

	/*
		this data is in msg.Envelope
		body and attachment is in msg.Body which is gotten by msg.GetBody()

	*/

	/*
		Every IMAP message (*imap.Message) that you fetch can include:
			the envelope → structured metadata (who sent it, subject, when, etc)
			the body → raw content of the email (MIME parts, attachments, etc)
	*/
	if msg.Envelope != nil {
		if len(msg.Envelope.From) > 0 {
			from = msg.Envelope.From[0].Address()
		}
		for _, t := range msg.Envelope.To {
			to = append(to, t.Address())
		}
		subject = msg.Envelope.Subject
		date = msg.Envelope.Date
	}

	// Append the parsed email to the slice
	emails = append(emails, models.Email{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
		Date:    date,
	})

	return emails
}

/*
GetBody gives you the raw message stream — everything including attachments, but not in a usable form yet.
That’s why we pass it into mail.CreateReader, which helps you separate:
but createReader also gives u stream but at least divide it into sections
inline parts (body text)

attachments (disposition: attachment)
*/
