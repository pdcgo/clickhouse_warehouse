-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hold_fund_history
(
    shop_id Int64,
    team_id Int64,
    day Date,

    hold_tx_count AggregateFunction(uniq, UInt64),
    hold_revenue_amount AggregateFunction(sum, Decimal(38, 18)),
    completed_tx_count AggregateFunction(uniq, UInt64),
    completed_revenue_amount AggregateFunction(sum, Decimal(38, 18))
)
ENGINE = AggregatingMergeTree
ORDER BY (shop_id, team_id, day)
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW hold_fund_history_mv
TO hold_fund_history
AS
SELECT
    shop_id,
    team_id,
    toDate(now('Asia/Jakarta')) AS day,
    uniqIfState(order_id,
        status NOT IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS hold_tx_count,
    sumIfState(revenue_amount,
        status NOT IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS hold_revenue_amount,
    uniqIfState(order_id,
        status IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS completed_tx_count,
    sumIfState(revenue_amount,
        status IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS completed_revenue_amount
FROM
(
    SELECT
        order_id,
        argMax(shop_id, _peerdb_version)        AS shop_id,
        argMax(team_id, _peerdb_version)        AS team_id,
        argMax(status, _peerdb_version)         AS status,
        argMax(revenue_amount, _peerdb_version) AS revenue_amount
    FROM order_latest
    GROUP BY order_id
)
GROUP BY 
(
    shop_id,
    team_id,
    day
)
-- +goose StatementEnd

-- +goose StatementBegin
-- backfilling query
INSERT INTO hold_fund_history
SELECT
    shop_id,
    team_id,
    toDate(now('Asia/Jakarta')) AS day,
    uniqIfState(order_id,
        status NOT IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS hold_tx_count,
    sumIfState(revenue_amount,
        status NOT IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS hold_revenue_amount,
    uniqIfState(order_id,
        status IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS completed_tx_count,
    sumIfState(revenue_amount,
        status IN (
            'completed',
            'cancel',
            'return_completed',
            'return_problem'
        )
    ) AS completed_revenue_amount
FROM
(
    SELECT
        order_id,
        argMax(shop_id, _peerdb_version)        AS shop_id,
        argMax(team_id, _peerdb_version)        AS team_id,
        argMax(status, _peerdb_version)         AS status,
        argMax(revenue_amount, _peerdb_version) AS revenue_amount
    FROM order_latest
    GROUP BY order_id
)
GROUP BY 
(
    shop_id,
    team_id,
    day
)
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS hold_fund_history_mv
-- +goose StatementEnd
-- +goose StatementBegin
-- DROP TABLE IF EXISTS hold_fund_history
-- +goose StatementEnd