-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    userid uuid PRIMARY KEY,
    username text NOT NULL,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS links (
    id varchar(8) PRIMARY KEY, 
    url text NOT NULL, 
    userid uuid NOT NULL REFERENCES users(uuid), 
    created_at timestamp NOT NULL DEFAULT NOW(),
    last_access timestamp NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS users;