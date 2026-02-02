-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders_latest
(
    `id` Int64,
    `team_id` Int64,
    `created_by_id` Int64,
    `status` String,
    `created_at` DateTime64(6),
    `order_mp_total` Int64,
    `order_mp_id` Int64,
    `_peerdb_version` UInt64
)
ENGINE = ReplacingMergeTree(_peerdb_version)
PRIMARY KEY id
ORDER BY id
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW orders_latest_mv
TO orders_latest
AS
SELECT * from orders
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW orders_latest_mv
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE orders_latest
-- +goose StatementEnd
