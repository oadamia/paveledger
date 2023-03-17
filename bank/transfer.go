package bank

import (
	"context"

	"encore.app/bank/model"
)

// CreateTransfer adds a new transfer
//
//encore:api public method=POST path=/transfers
func (s *Service) CreateTransfer(ctx context.Context, t *model.Transfer) error {
	err := validateTransfer(t)
	if err != nil {
		return err
	}

	return s.ledger.AddTransfer(ctx, t)
}

func validateTransfer(t *model.Transfer) error {
	return nil
}
