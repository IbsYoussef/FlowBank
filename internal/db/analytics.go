package db

import (
	"context"
	"flowbank/internal/model"
	"math"
)

// GetAnalytics retrieves aggregated analytics data for the dashboard
func (d *DB) GetAnalytics(ctx context.Context) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{}

	// 1. Get overall transaction statistics
	var avgAmount float64 // float64 for AVG result
	err := d.pool.QueryRow(ctx, `
		SELECT 
			COUNT(*) as total_transactions,
			COALESCE(SUM(ABS(amount)), 0) as total_volume,
			COALESCE(AVG(ABS(amount)), 0) as average_amount,
			COUNT(CASE WHEN transaction_type = 'credit' THEN 1 END) as credit_count,
			COUNT(CASE WHEN transaction_type = 'debit' THEN 1 END) as debit_count,
			COALESCE(SUM(CASE WHEN transaction_type = 'credit' THEN amount ELSE 0 END), 0) as money_in,
			COALESCE(SUM(CASE WHEN transaction_type = 'debit' THEN ABS(amount) ELSE 0 END), 0) as money_out
		FROM transactions
		WHERE status = 'COMPLETED'
	`).Scan(
		&stats.TotalTransactions,
		&stats.TotalVolume,
		&avgAmount, // scan into float64
		&stats.CreditCount,
		&stats.DebitCount,
		&stats.MoneyIn,
		&stats.MoneyOut,
	)
	if err != nil {
		return nil, err
	}

	// Convert average to int64
	stats.AverageAmount = int64(avgAmount)

	// 2. Get top merchants by transaction count
	merchantRows, err := d.pool.Query(ctx, `
		SELECT 
			merchant_name,
			COUNT(*) as count,
			SUM(ABS(amount)) as total_amount
		FROM transactions
		WHERE status = 'COMPLETED' AND merchant_name IS NOT NULL AND merchant_name != ''
		GROUP BY merchant_name
		ORDER BY count DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer merchantRows.Close()

	stats.TopMerchants = []model.MerchantStats{}
	for merchantRows.Next() {
		var m model.MerchantStats
		err := merchantRows.Scan(&m.MerchantName, &m.Count, &m.TotalAmount)
		if err != nil {
			return nil, err
		}
		stats.TopMerchants = append(stats.TopMerchants, m)
	}

	// 3. Get recent transactions (last 10)
	recentRows, err := d.pool.Query(ctx, `
		SELECT 
			t.transaction_id,
			u.user_name,
			t.amount,
			t.transaction_type,
			t.merchant_name,
			t.status,
			t.created_at
		FROM transactions t
		JOIN users u ON t.user_id = u.user_id
		ORDER BY t.created_at DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer recentRows.Close()

	stats.RecentTransactions = []model.TransactionDTO{}
	for recentRows.Next() {
		var tx model.TransactionDTO
		var amount int64
		var createdAt interface{}
		err := recentRows.Scan(
			&tx.ID,
			&tx.UserName,
			&amount,
			&tx.Type,
			&tx.MerchantName,
			&tx.Status,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		tx.Amount = int64(math.Abs(float64(amount)))
		tx.CreatedAt = createdAt.(interface{ String() string }).String()
		stats.RecentTransactions = append(stats.RecentTransactions, tx)
	}

	// 4. Get hourly volume (last 24 hours)
	hourlyRows, err := d.pool.Query(ctx, `
		SELECT 
			DATE_TRUNC('hour', created_at) as hour,
			SUM(CASE WHEN transaction_type = 'credit' THEN amount ELSE 0 END) as credit_total,
			SUM(CASE WHEN transaction_type = 'debit' THEN ABS(amount) ELSE 0 END) as debit_total,
			COUNT(*) as count
		FROM transactions
		WHERE created_at >= NOW() - INTERVAL '24 hours' AND status = 'COMPLETED'
		GROUP BY hour
		ORDER BY hour ASC
	`)
	if err != nil {
		return nil, err
	}
	defer hourlyRows.Close()

	stats.HourlyVolume = []model.HourlyVolumeData{}
	for hourlyRows.Next() {
		var hv model.HourlyVolumeData
		var hour interface{}
		err := hourlyRows.Scan(
			&hour,
			&hv.CreditTotal,
			&hv.DebitTotal,
			&hv.Count,
		)
		if err != nil {
			return nil, err
		}
		hv.Hour = hour.(interface{ String() string }).String()
		stats.HourlyVolume = append(stats.HourlyVolume, hv)
	}

	return stats, nil
}
