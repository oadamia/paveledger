package ledger

import (
	"context"

	"encore.app/ledger/bank"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const taskQueue = "authorization"

//encore:service
type Service struct {
	bank *bank.Bank

	tempClient client.Client
	tempWorker worker.Worker
}

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	b, err := bank.InitBank()
	if err != nil {
		return nil, err
	}

	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, err
	}

	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterActivity(b.CreateAccount)
	err = w.Start()
	if err != nil {
		c.Close()
		return nil, err
	}
	return &Service{bank: b, tempClient: c, tempWorker: w}, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.bank.Close()
	s.tempClient.Close()
	s.tempWorker.Stop()
}
