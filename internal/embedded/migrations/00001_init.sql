-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email text NOT NULL,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS links (
    id varchar(8) PRIMARY KEY, 
    url text NOT NULL, 
    user_uuid uuid NOT NULL REFERENCES users(uuid), 
    created_at timestamp NOT NULL DEFAULT NOW(),
    last_access timestamp NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";