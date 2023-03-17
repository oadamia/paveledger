package ledger

import (
	"context"
	"errors"

	"encore.app/bank/model"
	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

func (l *Ledger) AddTransfer(ctx context.Context, t *model.Transfer) error {
	res, err := l.client.CreateTransfers([]tb_types.Transfer{tbTransferFrom(*t)})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.TransferOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func (l *Ledger) AddPendingTransfer(ctx context.Context, t *model.Authorization) error {

	// res, err := l.client.CreateTransfers([]tb_types.Transfer{tbTransferFrom(*t)})
	// if err != nil {
	// 	return err
	// }

	// if len(res) > 0 && res[0].Result != tb_types.TransferOK {
	// 	return errors.New(res[0].Result.String())
	// }

	return nil
}

func tbTransferFrom(t model.Transfer) tb_types.Transfer {
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
