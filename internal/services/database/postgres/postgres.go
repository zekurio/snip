package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/embedded"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
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

func (p *Postgres) AddUpdateUser(user *models.User) error {
	query := `INSERT INTO users (uuid, email, password) VALUES ($1, $2, $3) ON CONFLICT (uuid) DO UPDATE SET email = $2, password = $3`
	_, err := p.db.Exec(query, user.UUID, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT uuid, password FROM users WHERE username = $1`
	var user models.User
	err := p.db.QueryRow(query, username).Scan(&user.UUID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email %s", username)
		}
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) GetUserByID(id string) (*models.User, error) {
	query := `SELECT email, password FROM users WHERE uuid = $1`
	var user models.User
	err := p.db.QueryRow(query, id).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %s", id)
		}
		return nil, err
	}

	return &user, nil
}

// Links

func (p *Postgres) AddUpdateLink(link *models.Link) error {
	query := `INSERT INTO links (id, url, user_uuid, created_at, last_access) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET url = $2, last_access = $5`
	_, err := p.db.Exec(query, link.ID, link.URL, link.UserID, link.CreatedAt, link.LastAccess)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetLinkByID(id string) (*models.Link, error) {
	query := `SELECT url, user_uuid, created_at, last_access FROM links WHERE id = $1`
	var link models.Link
	err := p.db.QueryRow(query, id).Scan(&link.URL, &link.UserID, &link.CreatedAt, &link.LastAccess)
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

		link.UserID = uuid
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
