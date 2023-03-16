package ledger

import (
	"context"
	"errors"

	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

// Transfer Representation
type TransferRepr struct {
	ID              string `json:"trasnfer_id"`
	DebitAccountID  string `json:"debit_account_id"`
	CreditAccountID string `json:"credit_account_id"`
	PendingID       string `json:"pending_id"`
	Ledger          uint32 `json:"ledger"`
	Code            uint16 `json:"code"`
	Amount          uint64 `json:"amount"`
	IsLinked        bool   `json:"is_linked"`
	IsPending       bool   `json:"is_pending_trasnfer"`
	IsPostPending   bool   `json:"is_post_pending_id"`
	IsVoidPending   bool   `json:"is_void_pending_id"`
}

// CreateTransfer adds a new transfer
//
//encore:api public method=POST path=/transfers
func (s *Service) CreateTransfer(ctx context.Context, t *TransferRepr) error {
	err := validateTransfer(t)
	if err != nil {
		return err
	}

	res, err := s.client.CreateTransfers([]tb_types.Transfer{transferFrom(*t)})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.TransferOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func validateTransfer(t *TransferRepr) error {
	return nil
}

func transferFrom(t TransferRepr) tb_types.Transfer {
	flag := tb_types.TransferFlags{
		Linked:              t.IsLinked,
		Pending:             t.IsPending,
		PostPendingTransfer: t.IsPostPending,
		VoidPendingTransfer: t.IsVoidPending,
	}

	return tb_types.Transfer{
		ID:              uint128(t.ID),
		DebitAccountID:  uint128(t.DebitAccountID),
		CreditAccountID: uint128(t.CreditAccountID),
		PendingID:       uint128(t.PendingID),
		Ledger:          t.Ledger,
		Code:            t.Code,
		Amount:          t.Amount,
		Flags:           flag.ToUint16(),
	}
}
