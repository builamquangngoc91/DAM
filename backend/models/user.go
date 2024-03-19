package models

import "time"

type User struct {
	UserID       string
	Username     string
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
