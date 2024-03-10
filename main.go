package main

import (
	"fmt"
	"groupieTrack/controller"
	"groupieTrack/manager"
	"groupieTrack/roots"
	initTemplate "groupieTrack/templates"

	"github.com/stripe/stripe-go"
)

func main() {
	stripe.Key = controller.StripeSecretKey
	manager.PrintColorResult("purple", "server is running...")
	fmt.Println("")
	manager.PrintColorResult("yellow", "CLICK HERE to OPEN  PAGE--->")
	fmt.Println("")
	manager.PrintColorResult("blue", " http://localhost:8080/connexion \n")
	fmt.Println("")
	manager.PrintColorResult("blue", " http://localhost:8080/home \n")
	fmt.Println("")
	manager.PrintColorResult("green", "TO STOP THE SERVER , PRESS  'ctrl+C' ")
	initTemplate.InitTemplate()
	roots.InitServe()
}
