package bank

import (
	"log"

	tb "github.com/tigerbeetledb/tigerbeetle-go"
)

type Bank struct {
	tbClient tb.Client
}

func InitBank() (*Bank, error) {
	tbc, err := tb.NewClient(0, []string{"3000"}, 1)
	if err != nil {
		log.Printf("Error creating client: %s", err)
		return nil, err
	}
	return &Bank{tbClient: tbc}, nil
}

func (b *Bank) Close() {
	b.tbClient.Close()
}
