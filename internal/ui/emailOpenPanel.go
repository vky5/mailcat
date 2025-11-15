package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/db/models"
)

// EmailOpenPanel displays the full content of a selected email
type EmailOpenPanel struct {
	textView *tview.TextView
	email    *models.Email
}

// NewEmailOpenPanel creates a styled panel for displaying email content
func NewEmailOpenPanel() *EmailOpenPanel {
	ep := &EmailOpenPanel{
		textView: tview.NewTextView(),
	}

	ep.textView.
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetBorder(true).
		SetTitle(" 📧 Email Content ").
		SetBackgroundColor(tcell.NewRGBColor(18, 30, 40)).SetBorderAttributes(tcell.AttrDim)

	ep.textView.SetTextColor(tcell.NewRGBColor(200, 200, 200))

	ep.showPlaceholder()

	ep.textView.SetFocusFunc(func() {
		ep.textView.SetBorderColor(tcell.NewRGBColor(0, 191, 255))
	})
	ep.textView.SetBlurFunc(func() {
		ep.textView.SetBorderColor(tcell.ColorNone).SetBorderAttributes(tcell.AttrDim)
	})

	return ep
}

// showPlaceholder displays a message when no email is selected
func (ep *EmailOpenPanel) showPlaceholder() {
	placeholder := `


                          			[#00BFFF::b]📬 No Email Selected[-:-:-]





                    [#B0B0B0]Select an email from the list to view its content here[-]


`
	ep.textView.SetText(placeholder)
}

// SetEmail updates the panel with email content
func (ep *EmailOpenPanel) SetEmail(email models.Email) {
	ep.email = &email
	ep.render()
}

// render displays the email content with rich formatting
func (ep *EmailOpenPanel) render() {
	if ep.email == nil {
		ep.showPlaceholder()
		return
	}

	var content strings.Builder

	// Header separator
	content.WriteString("\n[#00BFFF]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[-]\n\n")

	// Subject (bold & colored based on read)
	subjectColor := "#FFD700" // gold for unread
	if ep.email.Read {
		subjectColor = "#00CED1" // turquoise for read
	}
	content.WriteString(fmt.Sprintf("[%s::b]📨 %s[-:-:-]\n\n", subjectColor, ep.email.Subject))

	// From field
	content.WriteString(fmt.Sprintf("[#87CEEB::b]From:[-:-:-] [#B0C4DE]%s[-]\n", ep.email.From))

	// To field
	if strings.TrimSpace(ep.email.To) != "" {
		content.WriteString(fmt.Sprintf("[#87CEEB::b]To:[-:-:-]   [#B0C4DE]%s[-]\n", ep.email.To))
	}

	// Date field
	dateStr := ep.email.Date.Format("Monday, Jan 2, 2006 at 3:04 PM")
	content.WriteString(fmt.Sprintf("[#87CEEB::b]Date:[-:-:-] [#B0C4DE]%s[-]\n", dateStr))

	// Attachments if any
	if ep.email.Attachments != "" {
		attachments := strings.Split(ep.email.Attachments, ",")
		content.WriteString(fmt.Sprintf("\n[#FFA500::b]📎 Attachments (%d):[-:-:-]\n", len(attachments)))
		for i, att := range attachments {
			att = strings.TrimSpace(att)
			if att != "" {
				content.WriteString(fmt.Sprintf("   [#FFB366]%d. %s[-]\n", i+1, att))
			}
		}
	}

	// Body section
	body := strings.TrimSpace(ep.email.Body)
	content.WriteString("\n")
	if body == "" {
		content.WriteString("[#778899::i]No content[-:-:-]\n")
	} else {
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				content.WriteString("\n")
			} else if strings.HasPrefix(line, ">") {
				// Quoted lines
				content.WriteString(fmt.Sprintf("[#9370DB]%s[-]\n", line))
			} else {
				content.WriteString(fmt.Sprintf("[#E0E0E0]%s[-]\n", line))
			}
		}
	}

	// Footer hint
	content.WriteString("\n[#778899]Press [#32CD32]r[-] to reply  •  [#32CD32]f[-] to forward  •  [#32CD32]d[-] to delete[-]\n")

	ep.textView.SetText(content.String())
	ep.textView.ScrollToBeginning()
}

// Clear resets the panel to placeholder state
func (ep *EmailOpenPanel) Clear() {
	ep.email = nil
	ep.showPlaceholder()
}

// GetEmail returns the currently displayed email
func (ep *EmailOpenPanel) GetEmail() *models.Email {
	return ep.email
}

// Primitive returns the tview primitive
func (ep *EmailOpenPanel) Primitive() tview.Primitive {
	return ep.textView
}

// SetInputCapture allows setting custom key handlers
func (ep *EmailOpenPanel) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	ep.textView.SetInputCapture(capture)
}
