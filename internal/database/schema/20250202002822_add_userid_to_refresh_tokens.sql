-- +goose Up
ALTER TABLE refresh_tokens 
ADD user_id uuid NOT NULL; 

ALTER TABLE refresh_tokens
ADD CONSTRAINT fk_refresh_tokens_users
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- +goose Down
ALTER TABLE refresh_tokens
DROP column user_id;