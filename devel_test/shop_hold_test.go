package devel_test

import (
	"testing"

	"github.com/pdcgo/clickhouse_warehouse/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type ShopHoldHistory struct {
	Day              string
	ShopID           int
	TeamID           int
	TransactionCount int
	RevenueAmount    float64
	VersionAt        int64
}

func TestShopHold(t *testing.T) {
	db := database.NewLocalDatabase()
	defer db.Close()

	_, err := db.ExecContext(t.Context(), `
		INSERT INTO order_holds VALUES
		(1, 1, 1, 'created', 12000, 1, now()),
		(1, 1, 1, 'completed', 12000, -1, now()),
		(2, 1, 1, 'created', 1300, 1, now()),
		(3, 1, 1, 'created', 1300, 1, now())
	`)

	assert.Nil(t, err)

	rows, err := db.QueryContext(t.Context(), `
		select 
			day,
			shop_id,
			team_id,
			sumMerge(transaction_count),
			sumMerge(revenue_amount)
		from shop_hold_histories
		group by (
			day,
			shop_id,
			team_id
		);
	`)

	assert.Nil(t, err)

	var dd *gorm.DB
	dd.ScanRows()

	for rows.Next() {
		var hist ShopHoldHistory
		err := rows.Scan(&hist)
		assert.Nil(t, err)

		t.Log(hist.Day)
	}
}
