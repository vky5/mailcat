package commands

import (
	"strings"

	"github.com/vky5/mailcat/internal/db"
	"github.com/vky5/mailcat/internal/db/models"
)

type AddAccount struct {
	step     int
	email    string
	password string
	host     string
	port     string
	secure   bool
}

func NewAddAccount() *AddAccount {
	return &AddAccount{step: 0}
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

func (ac *AddAccount) HandleInput(input string, ctx Context) bool {
	switch ac.step {

	case 0:
		ac.email = input
		ac.step++
		ctx.ShowPlaceholder("Enter password:")
		return false

	case 1:
		ac.password = input
		ac.step++
		ctx.ShowPlaceholder("Enter IMAP host (e.g. imap.gmail.com):")
		return false

	case 2:
		ac.host = input
		ac.step++
		ctx.ShowPlaceholder("Enter IMAP port (e.g. 993):")
		return false

	case 3:
		ac.port = input
        ac.step++
        ctx.ShowPlaceholder("Use secure connection? (y/n):")
        return false

	case 4:
		lower := strings.ToLower(strings.TrimSpace(input))
		ac.secure = lower == "y" || lower == "yes" || lower == "true"

		account := models.Account{
			Email:     ac.email,
			Password:  ac.password,
			Host:      ac.host,
			Port:      ac.port,
			Secure:    ac.secure,
		}

		if err := db.DB.Create(&account).Error; err != nil {
			ctx.ShowMessage("Failed to create account: " + err.Error())
			return true
		}

		ctx.ShowMessage("Account added successfully!")
		ctx.ShowPlaceholder("")
		return true
	}

	return true
}
