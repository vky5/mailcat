package models

import "time"

type Email struct {
	ID      uint `gorm:"primaryKey"`
	From    string
	To      []string `gorm:"serializer:json"`
	Subject string
	Body    string
	Date    time.Time
	Read    bool
	Attachments string
}
