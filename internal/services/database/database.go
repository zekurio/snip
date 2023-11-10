package database

import (
	"time"

	"github.com/zekurio/snip/internal/models"
)

// IDatabase is an interface for our database drivers
type IDatabase interface {
	Close() error

	// Users

	// AddUpdateUser adds or updates a user
	AddUpdateUser(user *models.User) error

	// GetUserByUsername returns a user with the given email
	GetUserByUsername(username string) (*models.User, error)

	// GetUserByID returns a user with the given id
	GetUserByID(uuid string) (*models.User, error)

	// DeleteUser deletes a user with the given id
	DeleteUser(uuid string) error

	// Links

	// AddUpdateLink adds or updates a link
	AddUpdateLink(link *models.Link) error

	// GetLinkByID returns a link with the given id
	GetLinkByID(id string) (*models.Link, error)

	// GetLinksByUser returns all links for the given user
	GetLinksByUser(uuid string) ([]*models.Link, error)

	// DeleteLink deletes a link with the given id
	DeleteLink(id string) error

	// Refresh Tokens

	// SetUserRefreshToken sets a refresh token for the given user
	SetUserRefreshToken(ident, token string, expires time.Time) error

	// GetUserByRefreshToken returns a user with the given refresh tokens
	GetUserByRefreshToken(token string) (ident string, expires time.Time, err error)

	// RevokeUserRefreshToken revokes a refresh token for the given user
	RevokeUserRefreshToken(ident string) error
}
