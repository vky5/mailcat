package db

import "github.com/vky5/mailcat/internal/db/models"

func GetAccountByEmail(email string) (*models.Account, error) {
	var acc models.Account
	if err := DB.Where("email = ?", email).First(&acc).Error; err != nil {
		return nil, err
	}

	return &acc, nil
}