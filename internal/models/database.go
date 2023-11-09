package models

import "time"

type User struct {
	UUID     string
	Username string
	Password string
}

type Link struct {
	ID         string
	URL        string
	UserID     string
	CreatedAt  time.Time
	LastAccess time.Time
}
