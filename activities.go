package app

import "log"

type Activities struct {
}

func (a *Activities) CreatePayment() error {
	log.Println("Create payment...")
	return nil
}

func (a *Activities) SendAbandonedCartEmail(email string) error {
	log.Println("Send abandoned cart email to " + email)
	return nil
}
