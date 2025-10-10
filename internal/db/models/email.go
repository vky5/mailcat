package models

import "time"

type Email struct {
	ID      uint `gorm:"primaryKey"`
	From    string
	To      []string
	Subject string
	Body    string
	Date    time.Time
}
