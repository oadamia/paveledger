package ledger

import (
	"context"
	"errors"

	"encore.app/bank/model"
	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

func (l *Ledger) GetBalance(ctx context.Context, accountID string) (*model.Balance, error) {
	accounts, err := l.client.LookupAccounts([]tb_types.Uint128{uint128(accountID)})
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

	return &model.Balance{
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
