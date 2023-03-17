package bank

import (
	"context"
	"fmt"

	"encore.app/bank/model"
)

// CreatePresentment adds a new presentment
//
//encore:api public method=POST path=/presentments
func (s *Service) CreatePresentment(ctx context.Context, pre *model.Presentment) error {
	workflowID := fmt.Sprintf("%v_%v", pre.AccountID, pre.Amount)
	status, err := s.checkWorkFlowStatus(ctx, workflowID)
	if err != nil {
		return err
	}

	if status == WORKFLOW_STATUS_CLOSED || status == WORKFLOW_STATUS_NOT_FOUND {
		err := s.checkBalance(ctx, pre.AccountID, pre.Amount)
		if err != nil {
			return err
		}

		err = s.ledger.AddPresentTransfer(ctx, pre)
		if err != nil {
			return err
		}

	} else {
		err := s.client.SignalWorkflow(context.Background(), workflowID, "", PRESENTMENT_CHANNEL, pre)
		if err != nil {
			return err
		}
	}

	return nil
}
