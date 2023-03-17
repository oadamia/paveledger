package ledger

import (
	"context"
	"errors"
	"math/rand"

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

func (l *Ledger) AddPendingTransfer(ctx context.Context, a *model.Authorization) (string, error) {
	a.PendingID = generateTransactionID("P", 3)
	pt := pendingTransfer(*a)

	res, err := l.client.CreateTransfers([]tb_types.Transfer{pt})
	if err != nil {
		return a.PendingID, err
	}

	if len(res) > 0 && res[0].Result != tb_types.TransferOK {
		return a.PendingID, errors.New(res[0].Result.String())
	}

	return a.PendingID, nil
}

func pendingTransfer(a model.Authorization) tb_types.Transfer {
	flag := tb_types.TransferFlags{
		Linked:              false,
		Pending:             true,
		PostPendingTransfer: false,
		VoidPendingTransfer: false,
	}

	return tb_types.Transfer{
		ID:              uint128(generateTransactionID("T", 3)),
		DebitAccountID:  uint128(cardAccountID),
		CreditAccountID: uint128(a.AccountID),
		PendingID:       uint128(a.PendingID),
		Ledger:          defaultLedgerID,
		Code:            defaultCode,
		Amount:          a.Amount,
		Flags:           flag.ToUint16(),
	}
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

func generateTransactionID(prefix string, length int) string {
	randChars := make([]byte, length)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return prefix + string(randChars)
}
