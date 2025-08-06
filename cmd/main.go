package main

import billing "github.com/doddeeph/billing-engine/internal"

func main() {
	app := billing.NewBillingApp()
	app.Start()
}
