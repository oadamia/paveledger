package ledger

import (
	"context"
	"log"

	tb "github.com/tigerbeetledb/tigerbeetle-go"
)

//encore:service
type Service struct {
	client tb.Client
}

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	client, err := tb.NewClient(0, []string{"3000"}, 1)
	if err != nil {
		log.Printf("Error creating client: %s", err)
		return nil, err
	}

	return &Service{client: client}, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.client.Close()
}
