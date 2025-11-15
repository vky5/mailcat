package imap

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/vky5/mailcat/internal/db/models"
)


// stripHTML removes all HTML tags and replaces common HTML entities.
func stripHTML(html string) string {
	var out strings.Builder
	inTag := false

	for _, r := range html {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				out.WriteRune(r)
			}
		}
	}

	clean := out.String()

	replacer := strings.NewReplacer(
		"&nbsp;", " ",
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", "\"",
		"&#39;", "'",
	)

	return strings.TrimSpace(replacer.Replace(clean))
}


// looksLikeHTML does a simple check for HTML tags.
func looksLikeHTML(s string) bool {
	l := strings.ToLower(s)
	return strings.Contains(l, "<html") ||
		strings.Contains(l, "<body") ||
		strings.Contains(l, "<div") ||
		strings.Contains(l, "<span") ||
		strings.Contains(l, "<p") ||
		strings.Contains(l, "<table")
}


// parseMails converts an IMAP message into a models.Email value.
func parseMails(msg *imap.Message, emails []models.Email) []models.Email {
	if msg == nil {
		return emails
	}

	var e models.Email

	// Parse envelope fields
	if msg.Envelope != nil {
		e.Subject = msg.Envelope.Subject

		froms := make([]string, len(msg.Envelope.From))
		for i, a := range msg.Envelope.From {
			if a.PersonalName != "" {
				froms[i] = a.PersonalName + " <" + a.MailboxName + "@" + a.HostName + ">"
			} else {
				froms[i] = a.MailboxName + "@" + a.HostName
			}
		}
		e.From = strings.Join(froms, ", ")

		tos := make([]string, len(msg.Envelope.To))
		for i, a := range msg.Envelope.To {
			if a.PersonalName != "" {
				tos[i] = a.PersonalName + " <" + a.MailboxName + "@" + a.HostName + ">"
			} else {
				tos[i] = a.MailboxName + "@" + a.HostName
			}
		}
		e.To = strings.Join(tos, ", ")

		e.Date = msg.Envelope.Date
	}

	// Get full BODY[]
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		emails = append(emails, e)
		return emails
	}

	raw, err := io.ReadAll(r)
	if err != nil {
		emails = append(emails, e)
		return emails
	}

	mr, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		e.Body = stripHTML(string(raw))
		emails = append(emails, e)
		return emails
	}

	ct := mr.Header.Get("Content-Type")
	mediaType, params, _ := mime.ParseMediaType(ct)
	cte := mr.Header.Get("Content-Transfer-Encoding")

	// Handle multipart email
	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			bodyBytes, _ := io.ReadAll(mr.Body)
			text := decodeBytes(bodyBytes, cte)
			if looksLikeHTML(text) {
				text = stripHTML(text)
			}
			e.Body = text
			emails = append(emails, e)
			return emails
		}

		mpr := multipart.NewReader(mr.Body, boundary)

		var plain string
		var html string

		for {
			part, err := mpr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}

			partCT := part.Header.Get("Content-Type")
			partType, _, _ := mime.ParseMediaType(partCT)
			partCTE := part.Header.Get("Content-Transfer-Encoding")

			data, _ := io.ReadAll(part)
			decoded := decodeBytes(data, partCTE)

			if partType == "text/plain" && plain == "" {
				plain = decoded
			}

			if partType == "text/html" && html == "" {
				html = decoded
			}
		}

		if plain != "" {
			e.Body = strings.TrimSpace(plain)
		} else if html != "" {
			e.Body = stripHTML(html)
		}

		emails = append(emails, e)
		return emails
	}

	// Handle single-part email
	bodyBytes, _ := io.ReadAll(mr.Body)
	text := decodeBytes(bodyBytes, cte)

	if looksLikeHTML(text) {
		text = stripHTML(text)
	}

	e.Body = strings.TrimSpace(text)
	emails = append(emails, e)
	return emails
}


// decodeBytes decodes email body based on transfer encoding.
func decodeBytes(b []byte, cte string) string {
	cte = strings.ToLower(strings.TrimSpace(cte))

	switch cte {
	case "base64":
		dst := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
		n, err := base64.StdEncoding.Decode(dst, b)
		if err == nil {
			return string(dst[:n])
		}
		return string(b)

	case "quoted-printable":
		reader := quotedprintable.NewReader(bytes.NewReader(b))
		out, err := io.ReadAll(reader)
		if err == nil {
			return string(out)
		}
		return string(b)
	}

	return string(b)
}
