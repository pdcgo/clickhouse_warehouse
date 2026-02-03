-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_hold_histories
(
    day Date,
    shop_id Int64,
    team_id Int64,

    transaction_count AggregateFunction(sum, Int64),
    revenue_amount AggregateFunction(sum, Decimal(38, 18)),
    version_at Datetime(6)
)
ENGINE=AggregatingMergeTree
ORDER BY (day, shop_id, team_id)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_hold_histories
-- +goose StatementEnd
