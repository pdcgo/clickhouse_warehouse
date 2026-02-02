-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
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
ORDER BY id
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders
-- +goose StatementEnd
