package bank

import (
	"context"
	"errors"

	"encore.app/bank/model"
)

// GetAccount get account information from Database
//
//encore:api public path=/accounts/:accountID
func (s *Service) GetAccount(ctx context.Context, accountID string) (*model.Account, error) {
	return s.ledger.GetAccount(ctx, accountID)
}

// CreateAccount adds a new account to the bank.
//
//encore:api public method=POST path=/accounts
func (s *Service) CreateAccount(ctx context.Context, a *model.Account) error {
	err := validateAccount(*a)
	if err != nil {
		return err
	}

	return s.ledger.AddAccount(ctx, a)
}

func validateAccount(a model.Account) error {
	if a.IsCreditBalance == a.IsDebitBalance {
		return errors.New("credit and debit balance flag is equal")
	}
	return nil
}
