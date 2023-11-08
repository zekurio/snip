package database

import (
	"github.com/zekurio/snip/internal/models"
)

type DatabaseConfig struct {
	Type     string
	Postgres PostgresConfig
	// TODO add redis config
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// Database is an interface for our database drivers
type Database interface {

	// Users

	// CreateUser creates a new user with the given email and password
	CreateUser(email, password string) (*models.User, error)

	// GetUserByEmail returns a user with the given email
	GetUserByEmail(email string) (*models.User, error)

	// GetUserByID returns a user with the given id
	GetUserByID(uuid string) (*models.User, error)

	// Links

	// CreateLink creates a new link with the given url and user
	CreateLink(url string, user *models.User) (*models.Link, error)

	// GetLinkByID returns a link with the given id
	GetLinkByID(id string) (*models.Link, error)

	// GetLinksByUser returns all links for the given user
	GetLinksByUser(user *models.User) ([]*models.Link, error)

	// DeleteLink deletes a link with the given id
	DeleteLink(id string) error
}
