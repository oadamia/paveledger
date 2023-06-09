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

func (l *Ledger) AddPresentTransfer(ctx context.Context, pre *model.Presentment) error {
	res, err := l.client.CreateTransfers([]tb_types.Transfer{presentTransfer(*pre)})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.TransferOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func (l *Ledger) AddPendingTransfer(ctx context.Context, a *model.Authorization) (string, error) {
	pt := pendingTransfer(*a)

	res, err := l.client.CreateTransfers([]tb_types.Transfer{pt})
	if err != nil {
		return "", err
	}

	if len(res) > 0 && res[0].Result != tb_types.TransferOK {
		return a.PendingID, errors.New(res[0].Result.String())
	}

	return pt.ID.String(), nil
}

func (l *Ledger) AddPostPendingTransfer(ctx context.Context, p *model.Presentment) error {
	ppt := postPendingTransfer(*p)

	res, err := l.client.CreateTransfers([]tb_types.Transfer{ppt})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.TransferOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func (l *Ledger) AddVoidPendingTransfer(ctx context.Context, p *model.Authorization) error {
	ppt := voidPendingTransfer(*p)

	res, err := l.client.CreateTransfers([]tb_types.Transfer{ppt})
	if err != nil {
		return err
	}

	if len(res) > 0 && res[0].Result != tb_types.TransferOK {
		return errors.New(res[0].Result.String())
	}

	return nil
}

func presentTransfer(pre model.Presentment) tb_types.Transfer {
	flag := tb_types.TransferFlags{
		Linked:              false,
		Pending:             false,
		PostPendingTransfer: false,
		VoidPendingTransfer: false,
	}

	return tb_types.Transfer{
		ID:              uint128(generateID(3)),
		DebitAccountID:  uint128(pre.AccountID),
		CreditAccountID: uint128(cardAccountID),
		PendingID:       uint128("0"),
		Ledger:          defaultLedgerID,
		Code:            defaultCode,
		Amount:          pre.Amount,
		Flags:           flag.ToUint16(),
	}
}

func pendingTransfer(a model.Authorization) tb_types.Transfer {
	flag := tb_types.TransferFlags{
		Linked:              false,
		Pending:             true,
		PostPendingTransfer: false,
		VoidPendingTransfer: false,
	}

	return tb_types.Transfer{
		ID:              uint128(generateID(3)),
		DebitAccountID:  uint128(a.AccountID),
		CreditAccountID: uint128(cardAccountID),
		PendingID:       uint128("0"),
		Ledger:          defaultLedgerID,
		Code:            defaultCode,
		Amount:          a.Amount,
		Flags:           flag.ToUint16(),
	}
}

func postPendingTransfer(p model.Presentment) tb_types.Transfer {
	flag := tb_types.TransferFlags{
		Linked:              false,
		Pending:             false,
		PostPendingTransfer: true,
		VoidPendingTransfer: false,
	}

	return tb_types.Transfer{
		ID:              uint128(generateID(3)),
		DebitAccountID:  uint128(p.AccountID),
		CreditAccountID: uint128(cardAccountID),
		PendingID:       uint128(p.PendingID),
		Ledger:          defaultLedgerID,
		Code:            defaultCode,
		Amount:          p.Amount,
		Flags:           flag.ToUint16(),
	}
}

func voidPendingTransfer(a model.Authorization) tb_types.Transfer {
	flag := tb_types.TransferFlags{
		Linked:              false,
		Pending:             false,
		PostPendingTransfer: false,
		VoidPendingTransfer: true,
	}

	return tb_types.Transfer{
		ID:              uint128(generateID(3)),
		DebitAccountID:  uint128(a.AccountID),
		CreditAccountID: uint128(cardAccountID),
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
		ID:              uint128(generateID(3)),
		DebitAccountID:  uint128(t.DebitAccountID),
		CreditAccountID: uint128(t.CreditAccountID),
		PendingID:       uint128(t.PendingID),
		Ledger:          defaultLedgerID,
		Code:            defaultCode,
		Amount:          t.Amount,
		Flags:           flag.ToUint16(),
	}
}

func generateID(length int) string {
	randChars := make([]byte, length)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return string(randChars)
}
