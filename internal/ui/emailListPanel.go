package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/db/models"
	"github.com/vky5/mailcat/internal/logger"
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
	logger.Info("NewEmailListPanel: Creating new email list panel")
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
		logger.Info("EmailListPanel: Row selected:", row, "col:", col)
		// row 1,5,9... are subject rows (we render 4 rows per email)
		if row%4 != 1 {
			logger.Info("EmailListPanel: Non-subject row clicked, ignoring")
			return
		}
		idx := (row - 1) / 4
		logger.Info("EmailListPanel: Calculated email index:", idx)
		if idx >= 0 && idx < len(el.emails) && el.onSelect != nil {
			logger.Info("EmailListPanel: Calling onSelect for email:", el.emails[idx].Subject)
			el.onSelect(el.emails[idx])
		} else {
			logger.Info("EmailListPanel: Invalid index or no emails")
		}
	})

	// Single set of focus/blur handlers - just border color, no re-rendering
	el.table.SetFocusFunc(func() {
		logger.Info("EmailListPanel: Setting focus border color")
		el.table.SetBorderColor(tcell.NewRGBColor(0, 191, 255))
	})
	el.table.SetBlurFunc(func() {
		logger.Info("EmailListPanel: Removing focus border color")
		el.table.SetBorderColor(tcell.ColorNone).SetBorderAttributes(tcell.AttrDim)
	})

	logger.Info("NewEmailListPanel: Email list panel created successfully")
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
	logger.Info("SetEmails: Called with", len(emails), "emails")
	el.emails = emails
	logger.Info("SetEmails: Calling render()")
	el.render()

	// After rendering, move selection to the first subject row only if we have emails
	rowCount := el.table.GetRowCount()
	logger.Info("SetEmails: Table has", rowCount, "rows after render")
	
	if len(el.emails) > 0 && rowCount > 1 {
		// We have actual emails, select the first subject row
		logger.Info("SetEmails: Selecting row 1 (first email)")
		el.table.Select(1, 0)
	} else {
		// No emails, select the "No Emails" message (row 2)
		logger.Info("SetEmails: No emails, selecting empty state message")
		el.table.Select(2, 0)
	}
	logger.Info("SetEmails: Completed successfully")
}

// render displays emails with rich formatting
func (el *EmailListPanel) render() {
	logger.Info("render: Starting render with", len(el.emails), "emails")
	el.table.Clear()
	logger.Info("render: Table cleared")

	if len(el.emails) == 0 {
		logger.Info("render: No emails, showing empty state")
		cell := tview.NewTableCell("[::b][#00BFFF]📭 No Emails[-:-:-]").
			SetAlign(tview.AlignCenter).
			SetSelectable(true).  // Make selectable so arrow keys don't freeze
			SetBackgroundColor(tcell.NewRGBColor(18, 30, 40))
		el.table.SetCell(2, 0, cell)
		logger.Info("render: Empty state cell added")
		return
	}

	logger.Info("render: Rendering", len(el.emails), "emails")
	row := 1
	for i, e := range el.emails {
		logger.Info("render: Processing email", i, "-", e.Subject)
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

		// Row 4: Divider (non-selectable spacing between emails)
		dividerText := "[#2F4F4F]" + strings.Repeat("─", el.maxWidth) + "[-]"
		dividerCell := tview.NewTableCell(dividerText).
			SetSelectable(false).
			SetExpansion(1).
			SetBackgroundColor(bgColor)
		el.table.SetCell(row+3, 0, dividerCell)
		el.table.SetCell(row+3, 1, tview.NewTableCell("").SetBackgroundColor(bgColor).SetSelectable(false))

		row += 4
		logger.Info("render: Email", i, "rendered, next row:", row)
	}

	// Ensure first email subject is selected initially
	r, c := el.table.GetSelection()
	logger.Info("render: Current selection - row:", r, "col:", c)
	if el.table.GetRowCount() > 1 && r == 0 && c == 0 {
		logger.Info("render: Selecting first subject row (row 1)")
		el.table.Select(1, 0)
	}
	logger.Info("render: Render completed")
}

// Primitive returns the tview primitive
func (el *EmailListPanel) Primitive() tview.Primitive {
	return el.table
}


// SetLoading shows a temporary loading spinner instead of emails
func (el *EmailListPanel) SetLoading(loader tview.Primitive) {
	logger.Info("SetLoading: Setting loading state")
	el.table.Clear()
	el.table.SetCell(0, 0, tview.NewTableCell("").SetSelectable(false))
	el.table.SetCell(1, 0, tview.NewTableCell("").SetSelectable(false))
	el.table.SetCell(2, 0, tview.NewTableCell("").SetSelectable(false))

	// Center loader
	cell := tview.NewTableCell("").
		SetSelectable(false).
		SetAlign(tview.AlignCenter)
	el.table.SetCell(3, 0, cell)

	// Insert loader widget
	el.table.SetCell(4, 0,
		tview.NewTableCell("").
			SetSelectable(false).
			SetExpansion(1).
			SetAlign(tview.AlignCenter).
			SetReference(loader),
	)
	logger.Info("SetLoading: Loading state set successfully")
}