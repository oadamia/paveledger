package ledger

import (
	"log"

	tb "github.com/tigerbeetledb/tigerbeetle-go"
)

const (
	bankAccountID   = "3"
	cardAccountID   = "6"
	defaultLedgerID = 1
	defaultCode     = 1
)

type Ledger struct {
	client tb.Client
}

func Init() (*Ledger, error) {
	c, err := tb.NewClient(0, []string{"3000"}, 1)
	if err != nil {
		log.Printf("Error creating client: %s", err)
		return nil, err
	}
	return &Ledger{client: c}, nil
}

func (l *Ledger) Close() {
	l.client.Close()
}
