-- +goose Up
-- +goose StatementBegin
INSERT INTO shop_daily_metric
SELECT 
    shop_id,
    team_id,
    toDate(toTimeZone(created_at, 'Asia/Jakarta')) AS day,
    uniqState(order_id) as transaction_count,
    sumState(revenue_amount) as revenue_amount
FROM order_latest
GROUP BY (shop_id, team_id, toDate(toTimeZone(created_at, 'Asia/Jakarta')))
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
