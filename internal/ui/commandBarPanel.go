package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vky5/mailcat/internal/logger"
)

// Command represents a single executable command in the bar.
type Command struct {
	Name        string
	Description string
	Placeholder string
	Execute     func(input string, cb *CommandBar)
}

// CommandBar manages the input field, display area, and registry of commands.
type CommandBar struct {
	box      *tview.Flex
	input    *tview.InputField
	hintText *tview.TextView
	commands map[string]*Command
	nextFunc func(string)
	app      *tview.Application
}

// NewCommandBar creates and returns a new CommandBar component.
func NewCommandBar(app *tview.Application) *CommandBar {
	logger.Info("=== Creating CommandBar ===")
	cb := &CommandBar{
		app:      app,
		commands: make(map[string]*Command),
	}

	cb.input = tview.NewInputField().
		SetLabel("cmd> ").
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetFieldTextColor(tcell.ColorWhite)

	cb.hintText = tview.NewTextView().
		SetDynamicColors(true).
		SetText("Type : for command mode").
		SetTextAlign(tview.AlignLeft)

	cb.box = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(cb.hintText, 1, 0, false).
		AddItem(cb.input, 1, 0, true)

	return cb
}

// RegisterCommand adds a new command to the registry.
func (cb *CommandBar) RegisterCommand(cmd *Command) {
	logger.Info("Registering command: ", cmd.Name)
	cb.commands[cmd.Name] = cmd
}

// ShowPlaceholder updates the placeholder for next input.
func (cb *CommandBar) ShowPlaceholder(text string) {
	cb.input.SetPlaceholder(text)
	cb.input.SetText("")
}

// SetNextFunc sets the function to run after next Enter press.
func (cb *CommandBar) SetNextFunc(f func(string)) {
	cb.nextFunc = f
}

// ShowMessage displays a hint or output text.
func (cb *CommandBar) ShowMessage(msg string) {
	cb.hintText.SetText(msg)
}

// handleInput processes user input and triggers commands.
func (cb *CommandBar) handleInput(input string) {
	input = strings.TrimSpace(input)
	cb.input.SetText("")

	// Multi-step command flow
	if cb.nextFunc != nil {
		next := cb.nextFunc
		cb.nextFunc = nil
		next(input)
		return
	}

	// Execute registered commands
	if cmd, ok := cb.commands[input]; ok {
		cb.ShowMessage("Running: " + cmd.Description)
		cmd.Execute(input, cb)
	} else {
		cb.ShowMessage("[red]Unknown command: " + input)
		logger.Warn("Unknown command input: " + input)
	}
}

// GetPrimitive returns the root Flex to add into layouts.
func (cb *CommandBar) GetPrimitive() tview.Primitive {
	return cb.box
}

