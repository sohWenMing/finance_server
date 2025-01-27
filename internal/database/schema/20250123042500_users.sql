-- +goose Up
CREATE TABLE users (
    id uuid NOT NULL,
    is_admin BOOLEAN NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
