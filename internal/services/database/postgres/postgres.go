package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/embedded"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/pkg/randutils"
	"golang.org/x/crypto/bcrypt"
)

type Postgres struct {
	database.IDatabase

	db *sql.DB
}

func NewPostgres(c models.PostgresConfig) (pg *Postgres, err error) {
	pg = new(Postgres)

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Database, c.Username, c.Password)
	pg.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = pg.db.Ping()
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedded.Migrations)
	goose.SetDialect("postgres")
	goose.SetLogger(logrus.StandardLogger())
	err = goose.Up(pg.db, "migrations")
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

func (p *Postgres) CreateLink(url, uuid string) (*models.Link, error) {
	id, err := randutils.GetRandBase64Str(8)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO links (id, url, uuid) VALUES ($1, $2, $3) RETURNING created_at`
	var createdAt time.Time

	err = p.db.QueryRow(query, id, url, uuid).Scan(&createdAt)

	if err != nil {
		return nil, err
	}

	return &models.Link{ID: id, URL: url, UserUUID: uuid, CreatedAt: createdAt}, nil
}

func (p *Postgres) GetLinkByID(id string) (*models.Link, error) {
	query := `SELECT url, user_uuid, created_at, last_access FROM links WHERE id = $1`
	var link models.Link
	err := p.db.QueryRow(query, id).Scan(&link.URL, &link.UserUUID, &link.CreatedAt, &link.LastAccess)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no link found with id %s", id)
		}
		return nil, err
	}

	return &link, nil
}

func (p *Postgres) GetLinksByUser(uuid string) ([]*models.Link, error) {
	query := `SELECT id, url, created_at, last_access FROM links WHERE user_uuid = $1`
	rows, err := p.db.Query(query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var link models.Link
		err = rows.Scan(&link.ID, &link.URL, &link.CreatedAt, &link.LastAccess)
		if err != nil {
			return nil, err
		}

		link.UserUUID = uuid
		links = append(links, &link)
	}

	return links, nil
}

func (p *Postgres) DeleteLink(id string) error {
	query := `DELETE FROM links WHERE id = $1`
	_, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
