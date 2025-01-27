-- +goose Up
CREATE TABLE balance_sheets (
    id uuid NOT NULL,
    ticker TEXT NOT NULL,
    fiscal_date_ending TIMESTAMP NOT NULL,
    total_assets BIGINT NOT NULL,
    intangible_assets BIGINT NOT NULL,
    total_liabilities BIGINT NOT NULL,
    common_stock BIGINT NOT NULL,
    common_stock_shares_outstanding BIGINT NOT NULL,
    created_on TIMESTAMP NOT NULL,
    updated_on TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE balance_sheets;