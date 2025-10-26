package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/db/models"
)

// EmailListPanel displays emails for a folder/account.
type EmailListPanel struct {
	table    *tview.Table
	emails   []models.Email
	onSelect func(email models.Email)
	maxWidth int
}

// NewEmailListPanel creates a styled table for email list.
func NewEmailListPanel(onSelect func(email models.Email)) *EmailListPanel {
	el := &EmailListPanel{
		table:    tview.NewTable(),
		onSelect: onSelect,
		maxWidth: 60,
	}

	// Table styling with gradient-like background
	el.table.SetBorder(true).
		SetTitle(" 📬 Emails ").
		// SetBorderColor(tcell.NewRGBColor(0, 191, 255)).
		SetBackgroundColor(tcell.NewRGBColor(18, 30, 40)).SetBorderAttributes(tcell.AttrDim)

	el.table.SetSelectable(true, false)
	el.table.SetFixed(0, 2) // 2 columns for subject + date

	// Highlight style for selection - bright blue
	el.table.SetSelectedStyle(
		tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.NewRGBColor(0, 100, 150)).
			Bold(true),
	)

	// Only subject rows trigger selection
	el.table.SetSelectedFunc(func(row, col int) {
		if row%4 != 1 {
			return
		}
		idx := (row - 1) / 4
		if idx >= 0 && idx < len(el.emails) && el.onSelect != nil {
			el.onSelect(el.emails[idx])
		}
	})

	// Redraw on focus/blur for clean highlighting
	el.table.SetBlurFunc(func() { el.render() })
	el.table.SetFocusFunc(func() { el.render() })

	el.table.SetFocusFunc(func() {
		el.table.SetBorderColor(tcell.NewRGBColor(0, 191, 255))
	})
	el.table.SetBlurFunc(func() {
		el.table.SetBorderColor(tcell.ColorNone).SetBorderAttributes(tcell.AttrDim)
	})

	return el
}

// truncateString limits string to maxLen chars with ellipsis
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-1]) + "…"
}

// getPreviewText extracts first line of email body for preview
func getPreviewText(body string, maxLen int) string {
	body = strings.TrimSpace(body)
	lines := strings.Split(body, "\n")
	preview := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			preview = line
			break
		}
	}
	if preview == "" {
		preview = "No preview available"
	}
	return truncateString(preview, maxLen)
}

// SetEmails updates emails and re-renders
func (el *EmailListPanel) SetEmails(emails []models.Email) {
	el.emails = emails
	el.render()
}

// render displays emails with rich formatting
func (el *EmailListPanel) render() {
	el.table.Clear()

	if len(el.emails) == 0 {
		cell := tview.NewTableCell("[::b][#00BFFF]📭 No Emails[-:-:-]").
			SetAlign(tview.AlignCenter).
			SetSelectable(false).
			SetBackgroundColor(tcell.NewRGBColor(18, 30, 40))
		el.table.SetCell(2, 0, cell)
		return
	}

	row := 1
	for i, e := range el.emails {
		bgColor := tcell.NewRGBColor(18, 30, 40)
		if i%2 == 0 {
			bgColor = tcell.NewRGBColor(22, 35, 48)
		}

		// Read/unread styling
		envelope := "✉️"
		subjectColor := "#B0B0B0"
		fromColor := "#87CEEB"
		previewColor := "#778899"
		dateColor := "#00CED1"
		style := ""
		if !e.Read {
			envelope = "📧"
			subjectColor = "#FFD700"
			fromColor = "#00BFFF"
			previewColor = "#B0C4DE"
			dateColor = "#32CD32"
			style = "::b"
		}

		// Priority/importance indicator
		priorityIcon := ""
		if strings.Contains(strings.ToLower(e.Subject), "urgent") ||
			strings.Contains(strings.ToLower(e.Subject), "important") {
			priorityIcon = " [#FF4500]⚠️[-]"
		}

		// Attachment icon
		attachmentInfo := ""
		if e.Attachments != "" {
			attachmentInfo = " [#FFA500]📎[-]"
		}

		// Row 1: Subject + icons (left) + Date (right)
		subjectText := fmt.Sprintf("[%s%s]%s %s%s%s[-:-:-]",
			subjectColor,
			style,
			envelope,
			truncateString(e.Subject, el.maxWidth-20),
			priorityIcon,
			attachmentInfo,
		)
		dateText := fmt.Sprintf("[%s]%s[-]", dateColor, e.Date.Format("Jan 2"))

		subjectCell := tview.NewTableCell(subjectText).
			SetAlign(tview.AlignLeft).
			SetBackgroundColor(bgColor).
			SetExpansion(1).
			SetSelectable(true)

		dateCell := tview.NewTableCell(dateText).
			SetAlign(tview.AlignRight).
			SetBackgroundColor(bgColor).
			SetSelectable(false)

		el.table.SetCell(row, 0, subjectCell)
		el.table.SetCell(row, 1, dateCell)

		// Row 2: Sender
		senderText := fmt.Sprintf("[%s]👤 %s[-]", fromColor, truncateString(e.From, el.maxWidth-3))
		senderCell := tview.NewTableCell(senderText).
			SetAlign(tview.AlignLeft).
			SetBackgroundColor(bgColor).
			SetExpansion(1).
			SetSelectable(false)
		el.table.SetCell(row+1, 0, senderCell)
		el.table.SetCell(row+1, 1, tview.NewTableCell("").SetBackgroundColor(bgColor).SetSelectable(false))

		// Row 3: Preview
		previewText := fmt.Sprintf("   [%s]%s[-]", previewColor, getPreviewText(e.Body, el.maxWidth-3))
		previewCell := tview.NewTableCell(previewText).
			SetAlign(tview.AlignLeft).
			SetBackgroundColor(bgColor).
			SetExpansion(1).
			SetSelectable(false)
		el.table.SetCell(row+2, 0, previewCell)
		el.table.SetCell(row+2, 1, tview.NewTableCell("").SetBackgroundColor(bgColor).SetSelectable(false))

		// Row 4: Divider
		dividerText := "[#2F4F4F]" + strings.Repeat("─", el.maxWidth) + "[-]"
		dividerCell := tview.NewTableCell(dividerText).
			SetSelectable(false).
			SetExpansion(1).
			SetBackgroundColor(bgColor)
		el.table.SetCell(row+3, 0, dividerCell)
		el.table.SetCell(row+3, 1, tview.NewTableCell("").SetBackgroundColor(bgColor).SetSelectable(false))

		row += 4
	}

	// Ensure first email subject is selected initially
	r, c := el.table.GetSelection()
	if el.table.GetRowCount() > 1 && r == 0 && c == 0 {
		el.table.Select(1, 0)
	}
}

// Primitive returns the tview primitive
func (el *EmailListPanel) Primitive() tview.Primitive {
	return el.table
}
