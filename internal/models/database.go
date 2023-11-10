package models

import "time"

type User struct {
	ID       string
	Username string
	Password string
}

type Link struct {
	ID         string
	URL        string
	OwnerID    string
	CreatedAt  time.Time
	LastAccess time.Time
}
