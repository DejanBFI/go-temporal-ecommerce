package main

import (
	"log"

	"go-temporal-ecommerce/app"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	activities := new(app.Activities)

	w := worker.New(c, app.CartTaskQueue, worker.Options{})
	w.RegisterActivity(activities.CreatePayment)
	w.RegisterActivity(activities.SendAbandonedCartEmail)
	w.RegisterWorkflow(app.CartWorkflow)

	if err = w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
