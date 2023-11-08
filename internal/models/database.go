package models

import "time"

type User struct {
	UUID     string
	Email    string
	Password string
}

type Link struct {
	ID         string
	URL        string
	User       User
	CreatedAt  time.Time
	LastAccess time.Time
}
