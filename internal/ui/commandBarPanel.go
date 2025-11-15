package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/commands"
	"github.com/vky5/mailcat/internal/logger"
)

type CommandBar struct {
	box      *tview.Flex
	input    *tview.InputField
	hintText *tview.TextView
	registry map[string]commands.Command
	active   commands.Command
	app      *tview.Application
}

// NewCommandBar creates and returns a new CommandBar component.
func NewCommandBar(app *tview.Application) *CommandBar {
	blue := tcell.NewRGBColor(0, 102, 204) // input box blue

	cb := &CommandBar{
		app:      app,
		registry: make(map[string]commands.Command),
	}

	// INPUT FIELD (blue background)
	cb.input = tview.NewInputField().
		SetLabel("cmd> ").
		SetFieldBackgroundColor(blue).
		SetFieldTextColor(tcell.ColorWhite).
		SetPlaceholderStyle(
			tcell.StyleDefault.
				Foreground(tcell.ColorWhite). // bright readable placeholder
				Background(blue),
		).
		SetLabelColor(tcell.ColorLightBlue)

	// HINT BOX (dark background, readable text)
	cb.hintText = tview.NewTextView()
	cb.hintText.SetDynamicColors(true)
	cb.hintText.SetText("[#D0D0D0]Type : for command mode") // soft white/gray
	cb.hintText.SetTextAlign(tview.AlignLeft)
	cb.hintText.SetBackgroundColor(tcell.NewRGBColor(18, 30, 40)) // dark

	// Layout (keep dark background)
	cb.box = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(cb.hintText, 1, 0, false).
		AddItem(cb.input, 1, 0, true)

	cb.box.SetBackgroundColor(tcell.NewRGBColor(18, 30, 40))

	return cb
}

// Register adds a new command to the registry.
func (cb *CommandBar) Register(cmd commands.Command) {
	cb.registry[cmd.Name()] = cmd
}

// Context Interface implementation
func (cb *CommandBar) ShowPlaceholder(msg string) {
	cb.input.SetPlaceholder(msg)
}

func (cb *CommandBar) ShowMessage(msg string) {
	cb.hintText.SetText(msg)
}

// --- Main FSM handler ---
func (cb *CommandBar) handleInput(input string) {
	input = strings.TrimSpace(input)
	cb.input.SetText("")

	// If a command is already active...
	if cb.active != nil {
		done := cb.active.HandleInput(input, cb)
		if done {
			cb.active = nil
		}
		return
	}

	// Start new command
	cmd, ok := cb.registry[input]
	if !ok {
		cb.ShowMessage("[red]Unknown command: " + input)
		logger.Warn("Unknown command input: " + input)
		return
	}

	cb.active = cmd
	cb.ShowMessage("Running: " + cmd.Description())
	cmd.Begin(cb)
}

// GetPrimitive returns the root Flex to add into layouts.
func (cb *CommandBar) GetPrimitive() tview.Primitive {
	return cb.box
}
