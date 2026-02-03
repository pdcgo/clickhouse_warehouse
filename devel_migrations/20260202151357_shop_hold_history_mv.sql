-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW shop_hold_histories_mv
TO shop_hold_histories
AS
select 
    toDate(now()) as day,
	shop_id,
	team_id,
	sumState(toInt64(sign)) as transaction_count,
	sumState(revenue_amount * sign) as revenue_amount,
	now() as version_at
from order_holds
group by (day, shop_id, team_id, version_at)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW shop_hold_histories_mv
-- +goose StatementEnd
