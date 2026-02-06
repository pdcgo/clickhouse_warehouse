-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW order_incomplete_view AS
SELECT
    d.shop_id,
    uniq(d.order_id) as transaction_count,
    sum(d.revenue_amount) as revenue_amount
FROM (
    SELECT 
        o.order_id,
        o.shop_id,
        argMax(o.status, _peerdb_version) as status,
        argMax(o.revenue_amount, _peerdb_version) as revenue_amount
    FROM order_latest o
    WHERE 
        o.status IN ('completed', 'cancel', 'return_completed', 'return_problem')
    GROUP BY o.order_id, o.shop_id
) d
GROUP BY d.shop_id
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS order_incomplete_view;
-- +goose StatementEnd
