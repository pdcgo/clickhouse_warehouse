-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_latest
(
    order_id UInt64,
    shop_id Int64,
    team_id Int64, 
    status String,
    revenue_amount Decimal(38, 18),
    warehouse_fee_amount Decimal(38, 18),
    total_cost_amount Decimal(38, 18),
    order_time_at DateTime64(6),
    created_at DateTime64(6),
    _peerdb_version UInt64,

)
ENGINE = ReplacingMergeTree(_peerdb_version)
ORDER BY order_id
-- +goose StatementEnd




-- +goose StatementBegin
CREATE MATERIALIZED VIEW order_latest_mv
TO order_latest
AS
SELECT 
    id as order_id,
    order_mp_id as shop_id,
    team_id,
    status,
    toDecimal128(order_mp_total, 18) as revenue_amount,
    toDecimal128(warehouse_fee, 18) as warehouse_fee_amount,
    toDecimal128(total, 18) as total_cost_amount,
    order_time as order_time_at,
    created_at,
    _peerdb_version
FROM warehouse_prod.orders
-- +goose StatementEnd

-- +goose StatementBegin
-- backfilling
INSERT INTO order_latest
SELECT 
    id as order_id,
    order_mp_id as shop_id,
    team_id,
    status,
    toDecimal128(order_mp_total, 18) as revenue_amount,
    toDecimal128(warehouse_fee, 18) as warehouse_fee_amount,
    toDecimal128(total, 18) as total_cost_amount,
    order_time as order_time_at,
    created_at,
    _peerdb_version
FROM warehouse_prod.orders
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW order_latest_mv
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS order_latest
-- +goose StatementEnd
