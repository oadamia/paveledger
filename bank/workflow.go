package bank

import (
	"time"

	"encore.app/bank/ledger"
	"encore.app/bank/model"

	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
)

const AuthorizationTimeout = 100 * time.Second

func AutorizeWorkflow(ctx workflow.Context, list model.AuthorizationList) error {
	logger := workflow.GetLogger(ctx)

	authorizeChannel := workflow.GetSignalChannel(ctx, AUTHORIZATION_CHANNEL)
	presentmentChannel := workflow.GetSignalChannel(ctx, PRESENTMENT_CHANNEL)

	var l *ledger.Ledger

	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(authorizeChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var auth model.Authorization
			err := mapstructure.Decode(signal, &auth)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}
			logger.Info("add auth")
			list.Add(auth)
		})

		selector.AddReceive(presentmentChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var pre model.Presentment
			var auth *model.Authorization
			err := mapstructure.Decode(signal, &pre)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			}

			ctx = workflow.WithActivityOptions(ctx, ao)

			if !list.IsEmpty() {
				auth = list.Pop()
				pre.PendingID = auth.PendingID
			}

			err = workflow.ExecuteActivity(ctx, l.AddPostPendingTransfer, pre).Get(ctx, nil)
			if err != nil {
				logger.Error("Error adding pending transfer: %v", err)
				return
			}
		})

		if len(list.Items) > 0 {
			selector.AddFuture(workflow.NewTimer(ctx, authorizationTimeout), func(f workflow.Future) {

				auth := list.Pop()
				ao := workflow.ActivityOptions{
					StartToCloseTimeout: time.Minute,
				}

				ctx = workflow.WithActivityOptions(ctx, ao)

				err := workflow.ExecuteActivity(ctx, l.AddVoidPendingTransfer, auth).Get(ctx, nil)
				if err != nil {
					logger.Error("Error sending email %v", err)
					return
				}
			})
		}

		selector.Select(ctx)

		if list.IsEmpty() {
			break
		}
	}

	return nil
}
