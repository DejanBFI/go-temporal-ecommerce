package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go-temporal-ecommerce/app"

	"go.temporal.io/sdk/client"
)

func main() {
	ctx := context.Background()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	workflowID := "CART-" + fmt.Sprintf("%d", time.Now().Unix())

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: app.CartTaskQueue,
	}

	state := app.CartState{
		Items: make([]app.CartItem, 0),
		Email: "iamme@omg.lol",
	}

	we, err := c.ExecuteWorkflow(ctx, options, app.CartWorkflow, state)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	for i := 1; i <= 100; i++ {
		productID := rand.Int()%10000 + 1

		addToCartSignal := app.AddToCartSignal{
			Route: app.RouteTypes.ADD_TO_CART,
			Item: app.CartItem{
				ProductID: productID,
				Quantity:  1,
			},
		}
		err = c.SignalWorkflow(ctx, workflowID, we.GetRunID(), app.CartMessagesSignal, addToCartSignal)
		if err != nil {
			log.Fatalln("Unable to signal workflow", err)
		}
	}

	time.Sleep(65 * time.Second)

	resp, err := c.QueryWorkflow(ctx, workflowID, we.GetRunID(), app.QueryTypes.GET_CART)
	if err != nil {
		log.Fatalln("Unable to query workflow", err)
	}

	var result any
	if err = resp.Get(&result); err != nil {
		log.Fatalln("Unable to get query result", err)
	}

	log.Println("Workflow completed with state:", result)
}
