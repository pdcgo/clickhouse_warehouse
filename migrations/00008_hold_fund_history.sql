-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hold_fund_history
(
    shop_id Int64,
    team_id Int64,
    day Date,

    transaction_count Int64,
    revenue_amount Decimal(38, 18)
)
ENGINE=ReplacingMergeTree
ORDER BY (shop_id, team_id, day)
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW hold_fund_history_mv
TO hold_fund_history
AS
SELECT 
    shop_id,
    team_id,
    toDate(now('Asia/Jakarta')) AS day,
    uniq(order_id) as transaction_count,
    sum(revenue_amount) as revenue_amount
FROM order_latest
WHERE status NOT IN ('completed', 'cancel', 'return_completed', 'return_problem')
GROUP BY (shop_id, team_id, day)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS hold_fund_history_mv
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS hold_fund_history
-- +goose StatementEnd