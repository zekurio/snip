-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    user_id varchar(36) NOT NULL PRIMARY KEY,
    username text NOT NULL,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS links (
    id varchar(8) PRIMARY KEY,
    redirect_url text NOT NULL,
    owner_id varchar(36) NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT NOW(),
    last_access timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    ident varchar(36) NOT NULL PRIMARY KEY,
    token  TEXT NOT NULL DEFAULT '',
    expires TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down

DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS refresh_tokens;