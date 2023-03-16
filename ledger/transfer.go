package ledger

import (
	"context"

	"encore.app/ledger/bank"
)

// CreateTransfer adds a new transfer
//
//encore:api public method=POST path=/transfers
func (s *Service) CreateTransfer(ctx context.Context, t *bank.Transfer) error {
	err := validateTransfer(t)
	if err != nil {
		return err
	}

	return s.bank.CreateTransfer(ctx, t)
}

func validateTransfer(t *bank.Transfer) error {
	return nil
}
