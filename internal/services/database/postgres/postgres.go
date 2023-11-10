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
	query := `INSERT INTO users (user_id, username, password) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET username = $2, password = $3`
	_, err := p.db.Exec(query, user.ID, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT user_id, password FROM users WHERE username = $1`
	var user models.User
	err := p.db.QueryRow(query, username).Scan(&user.ID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with username %s", username)
		}
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) GetUserByID(userID string) (*models.User, error) {
	query := `SELECT username, password FROM users WHERE user_id = $1`
	var user models.User
	err := p.db.QueryRow(query, userID).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %s", userID)
		}
		return nil, err
	}

	return &user, nil
}

// Links

func (p *Postgres) AddUpdateLink(link *models.Link) error {
	query := `INSERT INTO links (id, redirect_url, owner_id, created_at, last_access) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET redirect_url = $2, last_access = $5`
	_, err := p.db.Exec(query, link.ID, link.URL, link.OwnerID, link.CreatedAt, link.LastAccess)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetLinkByID(linkID string) (*models.Link, error) {
	query := `SELECT redirect_url, owner_id, created_at, last_access FROM links WHERE id = $1`
	var link models.Link
	err := p.db.QueryRow(query, linkID).Scan(&link.URL, &link.OwnerID, &link.CreatedAt, &link.LastAccess)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no link found with id %s", linkID)
		}
		return nil, err
	}

	return &link, nil
}

func (p *Postgres) GetLinksByUser(userID string) ([]*models.Link, error) {
	query := `SELECT id, redirect_url, created_at, last_access FROM links WHERE owner_id = $1`
	rows, err := p.db.Query(query, userID)
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

		link.OwnerID = userID
		links = append(links, &link)
	}

	return links, nil
}

func (p *Postgres) DeleteLink(linkID string) error {
	query := `DELETE FROM links WHERE id = $1`
	_, err := p.db.Exec(query, linkID)
	if err != nil {
		return err
	}

	return nil
}

// Refresh Tokens

func (p *Postgres) SetUserRefreshToken(ident, token string, expires time.Time) error {
	query := `INSERT INTO refresh_tokens (ident, token, expires) VALUES ($1, $2, $3) ON CONFLICT (ident) DO UPDATE SET token = $2, expires = $3`
	_, err := p.db.Exec(query, ident, token, expires)
	return err
}

func (p *Postgres) GetUserByRefreshToken(token string) (ident string, expires time.Time, err error) {
	query := `SELECT ident, expires FROM refresh_tokens WHERE token = $1`
	row := p.db.QueryRow(query, token)
	err = row.Scan(&ident, &expires)
	return ident, expires, err
}

func (p *Postgres) RevokeUserRefreshToken(ident string) error {
	query := `DELETE FROM refresh_tokens WHERE ident = $1`
	_, err := p.db.Exec(query, ident)
	return err
}
