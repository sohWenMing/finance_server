-- +goose Up

ALTER TABLE balance_sheets
ADD CONSTRAINT unique_balance_sheet_id UNIQUE(id);

ALTER TABLE users
ADD CONSTRAINT unique_user_id UNIQUE(id);

ALTER TABLE refresh_tokens
ADD CONSTRAINT unique_refresh_token_id UNIQUE(id);

-- +goose Down
ALTER TABLE balance_sheets
DROP CONSTRAINT unique_balance_sheet_id;

ALTER TABLE users
DROP CONSTRAINT unique_user_id;

ALTER TABLE refresh_tokens
DROP CONSTRAINT unique_refresh_token_id;
