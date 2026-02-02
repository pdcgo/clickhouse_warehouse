-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_daily_metric
(
    shop_id Int64,
    team_id Int64,
    day Date,
    transaction_count Int64,
    revenue_amount Decimal(38, 18),
    version UInt64
)
ENGINE=ReplacingMergeTree(version)
ORDER BY (shop_id, team_id, day)
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_daily_metric
-- +goose StatementEnd
