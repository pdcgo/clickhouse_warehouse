-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_latest
(
    id Int64,
    user_id Int64,
    team_id Int64,
    shop_name String,
    mp_type String,
    uri String,
    shop_username String,
    _peerdb_version UInt64
)
ENGINE = ReplacingMergeTree(_peerdb_version)
ORDER BY id
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW shop_latest_mv
TO shop_latest
AS
SELECT
    id,
    user_id,
    team_id,
    mp_name as shop_name,
    mp_type,
    uri,
    mp_username as shop_username,
    _peerdb_version
FROM warehouse_prod.marketplaces
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO shop_latest
SELECT
    id,
    user_id,
    team_id,
    mp_name as shop_name,
    mp_type,
    uri,
    mp_username as shop_username,
    _peerdb_version
FROM warehouse_prod.marketplaces
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS shop_latest_mv
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS shop_latest
-- +goose StatementEnd
