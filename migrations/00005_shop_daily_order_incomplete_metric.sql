-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_daily_order_incomplete_metric
(
    shop_id Int64,
    team_id Int64,
    day Date,
    transaction_count AggregateFunction(uniq, UInt64),
    revenue_amount AggregateFunction(sum, Decimal(38, 18))
)
ENGINE=AggregatingMergeTree
ORDER BY (shop_id, team_id, day)

-- +goose StatementEnd
-- +goose StatementBegin
CREATE MATERIALIZED VIEW shop_daily_order_incomplete_metric_mv
TO shop_daily_order_incomplete_metric
AS
SELECT 
    shop_id,
    team_id,
    toDate(toTimeZone(created_at, 'Asia/Jakarta')) AS day,
    uniqState(order_id) as transaction_count,
    sumState(revenue_amount) as revenue_amount
FROM order_latest
WHERE status NOT IN ('completed', 'cancel', 'return_completed', 'return_problem')
GROUP BY (shop_id, team_id, toDate(toTimeZone(created_at, 'Asia/Jakarta')))
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW shop_daily_order_incomplete_metric_mv
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS shop_daily_order_incomplete_metric
-- +goose StatementEnd
