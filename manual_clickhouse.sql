CREATE TABLE prod.goose_db_version
(
    `id` UInt64,
    `version_id` Int64,
    `is_applied` UInt8,
    `tstamp` DateTime DEFAULT now()
)
ENGINE = SharedMergeTree('/clickhouse/tables/{uuid}/{shard}', '{replica}')
ORDER BY id
SETTINGS index_granularity = 8192

INSERT INTO prod.goose_db_version (id)
VALUES (DEFAULT);

CREATE TABLE prod.goose_db_version
(
    version_id Int64,
    is_applied UInt8,
    tstamp DateTime DEFAULT now()
)
ENGINE = MergeTree
ORDER BY version_id;
INSERT INTO prod.goose_db_version
    (version_id, is_applied)
VALUES
    (1, 1);