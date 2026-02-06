-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_holds
(
    order_id UInt64,
    shop_id Int64,
    team_id Int64, 
    created_by_id Int64,
    status String,
    revenue_amount Decimal(38, 18),
    sign Int8
)
ENGINE = CollapsingMergeTree(sign)
ORDER BY order_id
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW order_holds_mv
TO order_holds
AS
SELECT
    o.id as order_id,
    o.order_mp_id as shop_id,
    o.team_id,
    o.created_by_id,
    o.status,
    toDecimal128(order_mp_total, 18) as revenue_amount,
    CASE 
        WHEN status = 'created' THEN 1
        WHEN status IN ('completed', 'cancel', 'return_completed', 'return_problem') THEN -1
        ELSE 0
    END as sign
FROM warehouse_prod.orders o
    WHERE 
        o.status IN ('completed', 'cancel', 'return_completed', 'return_problem', 'created')
        AND o.is_partial != true
        AND o.is_order_fake != true
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW order_holds_mv
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE order_holds
-- +goose StatementEnd
