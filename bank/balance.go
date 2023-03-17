package bank

import (
	"context"

	"encore.app/bank/model"
)

// GetBalance get balance information from bank
//
//encore:api public path=/balances/:accountID
func (s *Service) GetBalance(ctx context.Context, accountID string) (*model.Balance, error) {
	return s.ledger.GetBalance(ctx, accountID)
}
