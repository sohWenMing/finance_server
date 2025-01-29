-- +goose Up
CREATE TABLE refresh_tokens (
    id uuid NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_on TIMESTAMP NOT NULL,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE refresh_tokens;