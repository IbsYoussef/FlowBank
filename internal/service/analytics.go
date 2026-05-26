package service

import (
	"context"
)

// DashboardResponse is the API response with formatted amounts
type DashboardResponse struct {
	TotalTransactions  int64                     `json:"total_transactions"`
	TotalVolume        Amount                    `json:"total_volume"`
	AverageAmount      Amount                    `json:"average_amount"`
	CreditCount        int64                     `json:"credit_count"`
	DebitCount         int64                     `json:"debit_count"`
	MoneyIn            Amount                    `json:"money_in"`
	MoneyOut           Amount                    `json:"money_out"`
	TopMerchants       []MerchantResponseDTO     `json:"top_merchants"`
	RecentTransactions []TransactionResponseDTO  `json:"recent_transactions"`
	HourlyVolume       []HourlyVolumeResponseDTO `json:"hourly_volume"`
}

// MerchantResponseDTO with formatted amounts
type MerchantResponseDTO struct {
	MerchantName string `json:"merchant_name"`
	Count        int64  `json:"count"`
	TotalAmount  Amount `json:"total_amount"`
}

// TransactionResponseDTO with formatted amounts
type TransactionResponseDTO struct {
	ID           string `json:"id"`
	UserName     string `json:"user_name"`
	Amount       Amount `json:"amount"`
	Type         string `json:"type"`
	MerchantName string `json:"merchant_name"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

// HourlyVolumeResponseDTO with formatted amounts
type HourlyVolumeResponseDTO struct {
	Hour        string `json:"hour"`
	CreditTotal Amount `json:"credit_total"`
	DebitTotal  Amount `json:"debit_total"`
	Count       int64  `json:"count"`
}

// GetDashboardStats retrieves aggregated analytics and formats amounts
func (s *Service) GetDashboardStats(ctx context.Context) (*DashboardResponse, error) {
	stats, err := s.db.GetAnalytics(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs with formatted amounts
	response := &DashboardResponse{
		TotalTransactions: stats.TotalTransactions,
		TotalVolume:       Amount(stats.TotalVolume),
		AverageAmount:     Amount(stats.AverageAmount),
		CreditCount:       stats.CreditCount,
		DebitCount:        stats.DebitCount,
		MoneyIn:           Amount(stats.MoneyIn),
		MoneyOut:          Amount(stats.MoneyOut),
	}

	// Convert merchants
	response.TopMerchants = make([]MerchantResponseDTO, len(stats.TopMerchants))
	for i, m := range stats.TopMerchants {
		response.TopMerchants[i] = MerchantResponseDTO{
			MerchantName: m.MerchantName,
			Count:        m.Count,
			TotalAmount:  Amount(m.TotalAmount),
		}
	}

	// Convert transactions
	response.RecentTransactions = make([]TransactionResponseDTO, len(stats.RecentTransactions))
	for i, tx := range stats.RecentTransactions {
		response.RecentTransactions[i] = TransactionResponseDTO{
			ID:           tx.ID,
			UserName:     tx.UserName,
			Amount:       Amount(tx.Amount),
			Type:         tx.Type,
			MerchantName: tx.MerchantName,
			Status:       tx.Status,
			CreatedAt:    tx.CreatedAt,
		}
	}

	// Convert hourly volume
	response.HourlyVolume = make([]HourlyVolumeResponseDTO, len(stats.HourlyVolume))
	for i, hv := range stats.HourlyVolume {
		response.HourlyVolume[i] = HourlyVolumeResponseDTO{
			Hour:        hv.Hour,
			CreditTotal: Amount(hv.CreditTotal),
			DebitTotal:  Amount(hv.DebitTotal),
			Count:       hv.Count,
		}
	}

	return response, nil
}

// FraudScoreResponseDTO with formatted amounts
type FraudScoreResponseDTO struct {
	TransactionID    string   `json:"transaction_id"`
	RiskScore        string   `json:"risk_score"`
	Status           string   `json:"status"`
	TriggeredRules   []string `json:"triggered_rules"`
	Confidence       float64  `json:"confidence"`
	ProcessingTimeMS float64  `json:"processing_time_ms"`
	ScoredAt         string   `json:"scored_at"`
	Amount           Amount   `json:"amount"`
	Type             string   `json:"type"`
	MerchantName     string   `json:"merchant_name"`
	UserName         string   `json:"user_name"`
}

// GetFraudScores retrieves recent fraud scores formatted for the dashboard
func (s *Service) GetFraudScores(ctx context.Context, limit int) ([]FraudScoreResponseDTO, error) {
	scores, err := s.db.GetFraudScores(ctx, limit)
	if err != nil {
		return nil, err
	}

	response := make([]FraudScoreResponseDTO, len(scores))
	for i, score := range scores {
		response[i] = FraudScoreResponseDTO{
			TransactionID:    score.TransactionID,
			RiskScore:        score.RiskScore,
			Status:           score.Status,
			TriggeredRules:   score.TriggeredRules,
			Confidence:       score.Confidence,
			ProcessingTimeMS: score.ProcessingTimeMS,
			ScoredAt:         score.ScoredAt,
			Amount:           Amount(score.Amount),
			Type:             score.Type,
			MerchantName:     score.MerchantName,
			UserName:         score.UserName,
		}
	}

	return response, nil
}
