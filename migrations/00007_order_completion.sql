-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_completion
(
    order_id UInt64,
    completion_at DateTime64(6),
    _peerdb_version UInt64
)
ENGINE = ReplacingMergeTree(_peerdb_version)
ORDER BY order_id
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW order_completion_mv
TO order_completion
AS
SELECT 
    order_id,
    timestamp as completion_at,
    _peerdb_version
FROM warehouse_prod.order_timestamps
WHERE 
    order_status IN ('completed', 'cancel', 'return_completed', 'return_problem')
-- +goose StatementEnd

-- +goose StatementBegin
-- backfilling
INSERT INTO order_completion
SELECT 
    order_id,
    timestamp as completion_at,
    _peerdb_version
FROM warehouse_prod.order_timestamps
WHERE 
    order_status IN ('completed', 'cancel', 'return_completed', 'return_problem')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW order_completion_mv;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS order_completion;
-- +goose StatementEnd
