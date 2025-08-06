package main

import "github.com/doddeeph/billing-engine/internal/billing"

func main() {
	app := billing.NewBillingApp()
	app.Start()
}
