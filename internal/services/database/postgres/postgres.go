package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/embedded"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"golang.org/x/crypto/bcrypt"
)

type Postgres struct {
	database.Database

	db *sql.DB
}

func NewPostgres(c database.PostgresConfig) (p *Postgres, err error) {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Database, c.Username, c.Password)
	p.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = p.db.Ping()
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedded.Migrations)
	goose.SetDialect("postgres")
	goose.SetLogger(logrus.StandardLogger())
	err = goose.Up(p.db, "migrations/postgres")
	if err != nil {
		return nil, err
	}

	return
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

// Users

func (p *Postgres) CreateUser(email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING uuid`
	var uuid string
	err = p.db.QueryRow(query, email, string(hashedPassword)).Scan(&uuid)
	if err != nil {
		return nil, err
	}

	return &models.User{UUID: uuid, Email: email}, nil
}

func (p *Postgres) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT uuid, password FROM users WHERE email = $1`
	var user models.User
	err := p.db.QueryRow(query, email).Scan(&user.UUID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email %s", email)
		}
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) GetUserByID(id string) (*models.User, error) {
	query := `SELECT email, password FROM users WHERE uuid = $1`
	var user models.User
	err := p.db.QueryRow(query, id).Scan(&user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %s", id)
		}
		return nil, err
	}

	return &user, nil
}

// Links
