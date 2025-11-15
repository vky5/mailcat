package ui

import "github.com/rivo/tview"

// Small centered loader text
func NewLoader(text string) *tview.TextView {
    tv := tview.NewTextView().
        SetTextAlign(tview.AlignCenter).
        SetText("[#00BFFF]⏳ " + text + " ...[-]").
        SetDynamicColors(true)
    return tv
}
