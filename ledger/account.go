package ledger

import (
	"context"
	"errors"

	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

// Account representation
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
func (s *Service) CreateAccount(ctx context.Context, a *AccountRepr) error {
	err := validateAccount(*a)
	if err != nil {
		return err
	}

	res, err := s.client.CreateAccounts([]tb_types.Account{accountFrom(*a)})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.AccountOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func validateAccount(a AccountRepr) error {
	if a.IsCreditBalance == a.IsDebitBalance {
		return errors.New("credit and debit balance flag is equal")
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

func accountFrom(a AccountRepr) tb_types.Account {
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
