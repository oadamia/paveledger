package bank

import (
	"context"
	"errors"

	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

type Account struct {
	ID              string `json:"account_id"`
	Ledger          uint32 `json:"ledger"`
	Code            uint16 `json:"code"`
	IsLinked        bool   `json:"is_linked"`
	IsDebitBalance  bool   `json:"is_debit_balance"`
	IsCreditBalance bool   `json:"is_credit_balance"`
}

func (b *Bank) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	tbAccounts, err := b.tbClient.LookupAccounts([]tb_types.Uint128{uint128(accountID)})
	if err != nil {
		return nil, err
	}

	if len(tbAccounts) == 0 {
		return nil, nil
	}

	return accountFrom(tbAccounts[0]), nil
}

func (b *Bank) CreateAccount(ctx context.Context, a *Account) error {
	res, err := b.tbClient.CreateAccounts([]tb_types.Account{tbAccountFrom(*a)})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.AccountOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func accountFrom(a tb_types.Account) *Account {
	return &Account{
		ID:              a.ID.String(),
		Ledger:          a.Ledger,
		Code:            a.Code,
		IsLinked:        check(a.Flags, isLinkedFlag()),
		IsDebitBalance:  check(a.Flags, isDebitBalanceFlag()),
		IsCreditBalance: check(a.Flags, isCreditBalanceFlag()),
	}
}

func tbAccountFrom(a Account) tb_types.Account {
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
