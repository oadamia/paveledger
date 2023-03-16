package ledger

import (
	"context"
	"errors"

	"encore.app/ledger/bank"
)

// GetAccount get account information from Database
//
//encore:api public path=/accounts/:accountID
func (s *Service) GetAccount(ctx context.Context, accountID string) (*bank.Account, error) {
	return s.bank.GetAccount(ctx, accountID)
}

// CreateAccount adds a new account to the bank.
//
//encore:api public method=POST path=/accounts
func (s *Service) CreateAccount(ctx context.Context, a *bank.Account) error {
	err := validateAccount(*a)
	if err != nil {
		return err
	}

	return s.bank.CreateAccount(ctx, a)
}

// GetBalance get balance information from bank
//
//encore:api public path=/balances/:accountID
func (s *Service) GetBalance(ctx context.Context, accountID string) (*bank.Balance, error) {
	return s.bank.GetBalance(ctx, accountID)
}

func validateAccount(a bank.Account) error {
	if a.IsCreditBalance == a.IsDebitBalance {
		return errors.New("credit and debit balance flag is equal")
	}
	return nil
}
