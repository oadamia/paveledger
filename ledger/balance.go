package ledger

import (
	"context"
	"errors"
	"log"

	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

type BalanceResponse struct {
	AccountID string `json:"account_id"`
	Available uint64 `json:"availabe"`
	Balance   uint64 `json:"balance"`
}

// GetBalance get balance information from Database
//
//encore:api public path=/balances/:accountID
func (s *Service) GetBalance(ctx context.Context, accountID string) (*BalanceResponse, error) {
	accounts, err := s.client.LookupAccounts([]tb_types.Uint128{uint128(accountID)})
	if err != nil {
		log.Printf("Could not fetch accounts: %s", err)
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, nil
	}

	balance, availableBalance, err := calculateBalance(accounts[0])
	if err != nil {
		log.Printf("Could not calculate balance: %s", err)
		return nil, err
	}

	return &BalanceResponse{
		AccountID: accounts[0].ID.String(),
		Available: availableBalance,
		Balance:   balance}, nil
}

func calculateBalance(account tb_types.Account) (balance, availableBalance uint64, err error) {

	if checkFlag(account.Flags, CreditsMustNotExceedDebitsFlag()) {
		availableBalance = account.DebitsPosted - account.CreditsPosted
		balance = availableBalance + account.DebitsPending
		return
	} else if checkFlag(account.Flags, DebitsMustNotExceedCreditsFlag()) {
		availableBalance = account.CreditsPosted - account.DebitsPosted
		balance = availableBalance + account.CreditsPending
		return
	}

	err = errors.New("account flag is not set")
	return
}
