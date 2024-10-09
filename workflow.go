package app

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
)

const abandonedCartTimeout time.Duration = 35 * time.Second

type (
	CartItem struct {
		ProductID int
		Quantity  int
	}

	CartState struct {
		Items                  []CartItem
		Email                  string
		SentAbandonedCartEmail bool
	}
)

func CartWorkflow(ctx workflow.Context, state CartState) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, QueryTypes.GET_CART, func(input []byte) (CartState, error) {
		return state, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	channel := workflow.GetSignalChannel(ctx, CartMessagesSignal)

	activities := new(Activities)
	for {
		// Create a new Selector on each iteration of the loop means Temporal will pick the first
		// event that occurs each time: either receiving a signal, or responding to the timer.
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(channel, func(c workflow.ReceiveChannel, _ bool) {
			var signal any
			c.Receive(ctx, &signal)

			var routeSignal RouteSignal
			err := mapstructure.Decode(signal, &routeSignal)
			if err != nil {
				logger.Error("Unable to decode signal.", "Error", err)
				return
			}

			switch routeSignal.Route {
			case RouteTypes.ADD_TO_CART:
				var message AddToCartSignal
				err := mapstructure.Decode(signal, &message)
				if err != nil {
					logger.Error("Invalid signal type %v", err)
				}
				state.AddToCart(message.Item)
			case RouteTypes.REMOVE_FROM_CART:
				var message RemoveFromCartSignal
				err := mapstructure.Decode(signal, &message)
				if err != nil {
					logger.Error("Invalid signal type %v", err)
				}
				state.RemoveFromCart(message.Item)
			default:
			}
		})

		if !state.SentAbandonedCartEmail && len(state.Items) > 0 {
			timer := workflow.NewTimer(ctx, abandonedCartTimeout)
			selector.AddFuture(timer, func(f workflow.Future) {
				state.SentAbandonedCartEmail = true
				ao := workflow.ActivityOptions{
					StartToCloseTimeout: 10 * time.Second,
				}

				ctx = workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(
					ctx,
					activities.SendAbandonedCartEmail,
					state.Email,
				).Get(ctx, nil)
				if err != nil {
					logger.Error("Error sending email %v", err)
					return
				}
			})
		}

		selector.Select(ctx)
	}
}

// @@@SNIPSTART temporal-ecommerce-add-and-remove
func (state *CartState) AddToCart(item CartItem) {
	for i := range state.Items {
		if state.Items[i].ProductID != item.ProductID {
			continue
		}

		state.Items[i].Quantity += item.Quantity
		return
	}

	state.Items = append(state.Items, item)
}

func (state *CartState) RemoveFromCart(item CartItem) {
	for i := range state.Items {
		if state.Items[i].ProductID != item.ProductID {
			continue
		}

		state.Items[i].Quantity -= item.Quantity
		if state.Items[i].Quantity <= 0 {
			state.Items = append(state.Items[:i], state.Items[i+1:]...)
		}
		break
	}
}

// @@@SNIPEND
