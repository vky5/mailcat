package commands

import "strings"

// HelpCommand lists all available commands
// it needs access to the command registry, so we can inject that in the constructor
type HelpCommand struct {
	registry map[string]Command
}


// this takes the registry from CommandBar and inside it store the 
func NewHelpCommand(reg map[string]Command) *HelpCommand{
	return &HelpCommand{registry: reg}
}

func (c *HelpCommand) Name() string {
	return  "!help"
}

func (c *HelpCommand) Description() string {
	return "Show all available commands"
}

func (c *HelpCommand) Begin(ctx Context){
	var b strings.Builder

	b.WriteString("[yellow]Available Commands:\n")

	for _, cmd := range c.registry {
		b.WriteString("  ")
		b.WriteString(cmd.Name())
		b.WriteString(" - ")
		b.WriteString(cmd.Description())
		b.WriteString("\n")
	}

	ctx.ShowMessage(b.String())
}

// Help in a single step command
// No extra input expected so handleinput finishes immediately 
func (c *HelpCommand) HandleInput(input string, ctx Context) bool {
	return  true
}