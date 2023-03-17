package ledger

import (
	"context"
	"errors"

	"encore.app/bank/model"
	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

func (l *Ledger) GetAccount(ctx context.Context, accountID string) (*model.Account, error) {
	tbAccounts, err := l.client.LookupAccounts([]tb_types.Uint128{uint128(accountID)})
	if err != nil {
		return nil, err
	}

	if len(tbAccounts) == 0 {
		return nil, nil
	}

	return accountFrom(tbAccounts[0]), nil
}

func (s *Ledger) AddAccount(ctx context.Context, a *model.Account) error {
	res, err := s.client.CreateAccounts([]tb_types.Account{tbAccountFrom(*a)})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.AccountOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func accountFrom(a tb_types.Account) *model.Account {
	return &model.Account{
		ID:              a.ID.String(),
		Ledger:          a.Ledger,
		Code:            a.Code,
		IsLinked:        check(a.Flags, isLinkedFlag()),
		IsDebitBalance:  check(a.Flags, isDebitBalanceFlag()),
		IsCreditBalance: check(a.Flags, isCreditBalanceFlag()),
	}
}

func tbAccountFrom(a model.Account) tb_types.Account {
	flag := tb_types.AccountFlags{
		Linked:                     a.IsLinked,
		DebitsMustNotExceedCredits: a.IsDebitBalance,
		CreditsMustNotExceedDebits: a.IsCreditBalance,
	}

	return tb_types.Account{
		ID:     uint128(a.ID),
		Ledger: a.Ledger,
		Code:   a.Code,
		Flags:  flag.ToUint16(),
	}
}
