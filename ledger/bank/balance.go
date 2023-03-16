package bank

import (
	"context"
	"errors"

	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

type Balance struct {
	AccountID        string `json:"account_id"`
	AvailableBalance uint64 `json:"availabe_balance"`
	AccountBalance   uint64 `json:"account_balance"`
}

func (b *Bank) GetBalance(ctx context.Context, accountID string) (*Balance, error) {
	accounts, err := b.tbClient.LookupAccounts([]tb_types.Uint128{uint128(accountID)})
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, nil
	}

	accountBalance, availableBalance, err := calculateBalance(accounts[0])
	if err != nil {
		return nil, err
	}

	return &Balance{
		AccountID:        accounts[0].ID.String(),
		AvailableBalance: availableBalance,
		AccountBalance:   accountBalance}, nil
}

func calculateBalance(account tb_types.Account) (accountBalance, availableBalance uint64, err error) {

	if check(account.Flags, isCreditBalanceFlag()) {
		availableBalance = account.DebitsPosted - account.CreditsPosted
		accountBalance = availableBalance + account.DebitsPending
		return
	} else if check(account.Flags, isDebitBalanceFlag()) {
		availableBalance = account.CreditsPosted - account.DebitsPosted
		accountBalance = availableBalance + account.CreditsPending
		return
	}

	err = errors.New("account flag is not set")
	return
}
