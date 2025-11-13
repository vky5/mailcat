package commands

// Command represents a state-machine command.
type Command interface {
	Name() string
	Description() string
	Begin (ctx Context)   // Called when the command begins.
	HandleInput(input string, ctx Context) (done bool)  // Called for each input until it returns done = true.
}


// Context is implemented by the UI layer.
// It gives commands the ability to display things WITHOUT knowing UI details
type Context interface {
	ShowMessage(text string)
	ShowPlaceholder(text string)
}

/*
using interface here so if any object that looks like a command 
satisfies the command interface

+ look at the synccommand it has no other staate simple func call it and ur work si done for syncing email
+ for addAccount there are different files and these files needs to be in the same 
- 

*/