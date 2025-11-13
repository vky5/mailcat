package commands

type AddAccount struct {
	step int
	email string
	Password string
}

func NewAddAccount(){

}

func (ac *AddAccount) Name() string {
	return "!addaccount"
}

func (ac *AddAccount) Description() string {
	return "Add IMAP account to the client"
}

func (ac *AddAccount) Begin(ctx Context) {
	ctx.ShowPlaceholder("Enter email:")
}

func (c *AddAccount) HandleInput(input string, ctx Context) bool {
	switch c.step {
	case 0:
		c.email = input
		c.step++
		ctx.ShowPlaceholder("Enter password:")
		return false

	case 1:
		password := input
		// TODO: add saving logic
		_ = password

		ctx.ShowMessage("Account added: " + c.email)
		return true
	}

	return true
}