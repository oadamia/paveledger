package bank

import (
	"context"

	"encore.app/bank/ledger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const taskQueue = "authorization"

//encore:service
type Service struct {
	ledger     *ledger.Ledger
	tempClient client.Client
	tempWorker worker.Worker
}

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	s, err := ledger.Init()
	if err != nil {
		return nil, err
	}

	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, err
	}

	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterActivity(s.AddAccount)
	err = w.Start()
	if err != nil {
		c.Close()
		return nil, err
	}
	return &Service{ledger: s, tempClient: c, tempWorker: w}, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.ledger.Close()
	s.tempClient.Close()
	s.tempWorker.Stop()
}
