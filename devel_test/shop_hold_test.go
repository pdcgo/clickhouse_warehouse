package devel_test

import (
	"testing"
	"time"

	"github.com/pdcgo/clickhouse_warehouse/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type ShopHoldHistory struct {
	Day              time.Time       `ch:"day"`
	ShopID           int64           `ch:"shop_id"`
	TeamID           int64           `ch:"team_id"`
	TransactionCount int64           `ch:"transaction_count"`
	RevenueAmount    decimal.Decimal `ch:"revenue_amount"`
	// VersionAt        time.Time `ch:"version_at"`
}

func TestShopHold(t *testing.T) {
	db, err := database.NewLocalDatabaseHttp()
	assert.Nil(t, err)

	defer db.Close()

	assert.Nil(t, err)

	err = db.Exec(t.Context(), `
		INSERT INTO order_holds VALUES
		(1, 1, 1, 'created', 12000, 1, now()),
		(1, 1, 1, 'completed', 12000, -1, now()),
		(2, 1, 1, 'created', 1300, 1, now()),
		(3, 1, 1, 'created', 1300, 1, now())
	`)

	assert.Nil(t, err)

	rows, err := db.Query(t.Context(), `
		select 
			day,
			shop_id,
			team_id,
			sumMerge(transaction_count) as transaction_count,
			sumMerge(revenue_amount) as revenue_amount
		from shop_hold_histories
		group by (
			day,
			shop_id,
			team_id
		);
	`)

	assert.Nil(t, err)

	for rows.Next() {
		var hist ShopHoldHistory
		err := rows.ScanStruct(&hist)
		assert.Nil(t, err)

		switch hist.ShopID {
		case 1:
			assert.Equal(t, int64(2), hist.TransactionCount)
			revAmount, _ := hist.RevenueAmount.Float64()
			assert.Equal(t, float64(2600), revAmount)
		}

		t.Log(hist.Day)
	}
}
