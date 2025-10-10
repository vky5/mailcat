package models

import "time"


// account credentials for IMAP
type Account struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex"`
	Password  string
	Secure    bool
	Host      string
	Port      string
	CreatedAt time.Time
}
