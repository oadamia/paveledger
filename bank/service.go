package bank

import (
	"context"

	"encore.app/bank/ledger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

//encore:service
type Service struct {
	ledger *ledger.Ledger
	client client.Client
	worker worker.Worker
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

	w := worker.New(c, TASK_QUEUE, worker.Options{})
	w.RegisterWorkflow(AutorizeWorkflow)
	w.RegisterActivity(s.AddPostPendingTransfer)
	w.RegisterActivity(s.AddVoidPendingTransfer)

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, err
	}
	return &Service{ledger: s, client: c, worker: w}, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.ledger.Close()
	s.client.Close()
	s.worker.Stop()
}
