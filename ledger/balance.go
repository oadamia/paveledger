package ledger

import (
	"context"
)

// PingResponse is the response from the Ping endpoint.
type BalanceResponse struct {
	Available int `json:"availabe"`
	Balance   int `json:"balance"`
}

// Ping pings a specific site and determines whether it's up or down right now.
//
//encore:api public path=/balances/:accountID
func (s *Service) Balances(ctx context.Context, accountID int) (*BalanceResponse, error) {
	return &BalanceResponse{Available: 100, Balance: 500}, nil
}