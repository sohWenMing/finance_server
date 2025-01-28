-- name: CreateBalanceSheetRecord :one
INSERT INTO balance_sheets (
id, 
ticker, 
fiscal_date_ending, 
total_assets, 
intangible_assets, 
total_liabilities,
common_stock,
common_stock_shares_outstanding,
created_on,
updated_on
)
VALUES($1,$2,$3,$4,$5,$6,$7,$8, $9, $10)
RETURNING *;


-- name: GetBalanceSheetRecordsByTicker :many
SELECT * from balance_sheets
WHERE ticker = $1;

-- name: DeleteBalanceSheetRecordsByTicker :exec
DELETE from balance_sheets
WHERE ticker = $1;