-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_holds
(
    order_id Int64,
    shop_id Int64,
    team_id Int64,

    status String,
    revenue_amount Decimal(38, 18),

    sign Int8,
    version_at Datetime(6)
)
ENGINE=CollapsingMergeTree(sign)
ORDER BY order_id
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_holds
-- +goose StatementEnd
