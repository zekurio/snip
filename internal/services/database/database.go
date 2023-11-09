package database

import (
	"github.com/zekurio/snip/internal/models"
)

// IDatabase is an interface for our database drivers
type IDatabase interface {
	Close() error

	// Users

	// CreateUser creates a new user with the given email and password
	CreateUser(email, password string) (*models.User, error)

	// GetUserByEmail returns a user with the given email
	GetUserByEmail(email string) (*models.User, error)

	// GetUserByID returns a user with the given id
	GetUserByID(uuid string) (*models.User, error)

	// Links

	// CreateLink creates a new link with the given url and user
	CreateLink(url, uuid string) (*models.Link, error)

	// GetLinkByID returns a link with the given id
	GetLinkByID(id string) (*models.Link, error)

	// GetLinksByUser returns all links for the given user
	GetLinksByUser(uuid string) ([]*models.Link, error)

	// DeleteLink deletes a link with the given id
	DeleteLink(id string) error
}
