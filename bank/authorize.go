package bank

import (
	"context"
	"errors"
	"fmt"
	"time"

	"encore.app/bank/model"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

// CreateAuthorization adds a new authorization
//
//encore:api public method=POST path=/authorizations
func (s *Service) CreateAuthorization(ctx context.Context, auth *model.Authorization) error {
	err := s.checkBalance(ctx, auth.AccountID, auth.Amount)
	if err != nil {
		return err
	}

	auth.PendingID, err = s.ledger.AddPendingTransfer(ctx, auth)
	if err != nil {
		return err
	}

	workflowID := fmt.Sprintf("%v_%v", auth.AccountID, auth.Amount)
	status, err := s.checkWorkFlowStatus(ctx, workflowID)
	if err != nil {
		return err
	}

	if status == WORKFLOW_STATUS_CLOSED || status == WORKFLOW_STATUS_NOT_FOUND {
		err := s.startWorkFlow(workflowID)
		if err != nil {
			return err
		}
	}

	auth.Timestamp = time.Now().Unix()
	return s.client.SignalWorkflow(context.Background(), workflowID, "", AUTHORIZATION_CHANNEL, auth)
}

func (s *Service) checkBalance(ctx context.Context, ID string, amount uint64) error {
	balance, err := s.ledger.GetBalance(ctx, ID)
	if err != nil {
		return err
	}

	if amount > balance.AvailableBalance {
		return errors.New("insufficient funds")
	}

	return nil
}

func (s *Service) checkWorkFlowStatus(ctx context.Context, workflowID string) (workflowStatus, error) {
	var notFoundError *serviceerror.NotFound
	res, err := s.client.DescribeWorkflowExecution(ctx, workflowID, "")

	if err != nil {
		if errors.As(err, &notFoundError) {
			return WORKFLOW_STATUS_NOT_FOUND, nil
		}

		return WORKFLOW_STATUS_NOT_FOUND, err
	}

	if res.WorkflowExecutionInfo.Status == enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
		return WORKFLOW_STATUS_RUNNING, nil
	} else {
		return WORKFLOW_STATUS_CLOSED, nil
	}
}

func (s *Service) startWorkFlow(workflowID string) error {
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: TASK_QUEUE,
	}

	authList := model.AuthorizationList{Items: make([]model.Authorization, 0)}
	_, err := s.client.ExecuteWorkflow(context.Background(), options, AutorizeWorkflow, authList)
	return err
}
