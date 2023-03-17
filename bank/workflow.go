package bank

import (
	"time"

	"encore.app/bank/ledger"
	"encore.app/bank/model"

	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
)

type authList struct {
	Items []model.Authorization
}

func (l *authList) Add(item model.Authorization) {
	l.Items = append(l.Items, item)
}

func (l *authList) Pop() *model.Authorization {
	if len(l.Items) == 0 {
		return nil
	}

	a := l.Items[0]
	l.Items = l.Items[1:]
	return &a
}

const AuthorizationTimeout = 100 * time.Second

func AutorizeWorkflow(ctx workflow.Context, alist authList) error {
	logger := workflow.GetLogger(ctx)
	err := workflow.SetQueryHandler(ctx, "authorize", func(input []byte) (authList, error) {
		return alist, nil
	})
	if err != nil {
		return err
	}

	authorizeChannel := workflow.GetSignalChannel(ctx, model.SignalChannels.AUTHORIZATION_CHANNEL)
	presentmentChannel := workflow.GetSignalChannel(ctx, model.SignalChannels.PRESENTMENT_CHANNEL)

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

			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			}

			ctx = workflow.WithActivityOptions(ctx, ao)

			err = workflow.ExecuteActivity(ctx, l.AddPendingTransfer, auth).Get(ctx, nil)
			if err != nil {
				logger.Error("Error creating stripe charge: %v", err)
				return
			}

			alist.Add(auth)
		})

		selector.AddReceive(presentmentChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var pre model.Presentment
			err := mapstructure.Decode(signal, &pre)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}
		})

		// selector.AddReceive(checkoutChannel, func(c workflow.ReceiveChannel, _ bool) {
		// 	var signal interface{}
		// 	c.Receive(ctx, &signal)

		// 	var message CheckoutSignal
		// 	err := mapstructure.Decode(signal, &message)
		// 	if err != nil {
		// 		logger.Error("Invalid signal type %v", err)
		// 		return
		// 	}

		// 	alist.Email = message.Email

		// 	ao := workflow.ActivityOptions{
		// 		StartToCloseTimeout: time.Minute,
		// 	}

		// 	ctx = workflow.WithActivityOptions(ctx, ao)

		// 	err = workflow.ExecuteActivity(ctx, a.CreateStripeCharge, alist).Get(ctx, nil)
		// 	if err != nil {
		// 		logger.Error("Error creating stripe charge: %v", err)
		// 		return
		// 	}

		// 	checkedOut = true
		// })

		// if !sentAbandonedCartEmail && len(alist.Items) > 0 {
		// 	selector.AddFuture(workflow.NewTimer(ctx, abandonedCartTimeout), func(f workflow.Future) {
		// 		sentAbandonedCartEmail = true
		// 		ao := workflow.ActivityOptions{
		// 			StartToCloseTimeout: time.Minute,
		// 		}

		// 		ctx = workflow.WithActivityOptions(ctx, ao)

		// 		err := workflow.ExecuteActivity(ctx, a.SendAbandonedCartEmail, alist.Email).Get(ctx, nil)
		// 		if err != nil {
		// 			logger.Error("Error sending email %v", err)
		// 			return
		// 		}
		// 	})
		// }

		// selector.Select(ctx)

		// if checkedOut {
		// 	break
		// }
	}

	return nil
}
