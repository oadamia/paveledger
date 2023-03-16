package ledger

import (
	"context"
	"errors"

	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

type AccountRepr struct {
	ID              string `json:"account_id"`
	Ledger          uint32 `json:"ledger"`
	Code            uint16 `json:"code"`
	IsLinked        bool   `json:"is_linked"`
	IsDebitBalance  bool   `json:"is_debit_balance"`
	IsCreditBalance bool   `json:"is_credit_balance"`
}

// GetAccount get account information from Database
//
//encore:api public path=/accounts/:accountID
func (s *Service) GetAccount(ctx context.Context, accountID string) (*AccountRepr, error) {
	accounts, err := s.client.LookupAccounts([]tb_types.Uint128{uint128(accountID)})
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, nil
	}

	return accountReprFrom(accounts[0]), nil

}

// CreateAccount adds a new account to the list of accounts.
//
//encore:api public method=POST path=/accounts
func (s *Service) CreateAccount(ctx context.Context, cap *AccountRepr) error {
	if cap.IsCreditBalance == cap.IsDebitBalance {
		return errors.New("credit and debit balance flag is equal")
	}

	res, err := s.client.CreateAccounts([]tb_types.Account{accountFrom(*cap)})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.AccountOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func accountReprFrom(a tb_types.Account) *AccountRepr {
	return &AccountRepr{
		ID:              a.ID.String(),
		Ledger:          a.Ledger,
		Code:            a.Code,
		IsLinked:        check(a.Flags, isLinkedFlag()),
		IsDebitBalance:  check(a.Flags, isDebitBalanceFlag()),
		IsCreditBalance: check(a.Flags, isCreditBalanceFlag()),
	}
}

func accountFrom(cap AccountRepr) tb_types.Account {
	flag := tb_types.AccountFlags{
		Linked:                     cap.IsLinked,
		DebitsMustNotExceedCredits: cap.IsDebitBalance,
		CreditsMustNotExceedDebits: cap.IsCreditBalance,
	}

	return tb_types.Account{
		ID:     uint128(cap.ID),
		Ledger: cap.Ledger,
		Code:   cap.Code,
		Flags:  flag.ToUint16(),
	}
}
