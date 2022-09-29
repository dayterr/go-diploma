package main

import (
	"github.com/dayterr/go-diploma/cmd/gophermart/handlers"
	"github.com/dayterr/go-diploma/internal/accrual"
	config2 "github.com/dayterr/go-diploma/internal/config"
	"log"
	"net/http"
)

func main() {
	config, err := config2.GetConfig()
	if err != nil {
		log.Fatal("no config, can't start the program")
	}

	orderChannel := make(chan string)
	ah := handlers.NewAsyncHandler(config.DatabaseURI)
	r := handlers.CreateRouterWithAsyncHandler(ah)

	ac := accrual.NewAccrualClient(config.AccrualSystemAddress, config.DatabaseURI, orderChannel)
	ah.AccrualClient = ac

	go func() {
		for {
			select {
			case <- ah.AccrualClient.OrderChannel:
				ah.AccrualClient.ManagePoints(<- ah.AccrualClient.OrderChannel)
			}
		}
	}()
	//go ah.AccrualClient.ReadOrderNumber()

	err = http.ListenAndServe(config.RunAddress, r)

	if err != nil {
		log.Fatal("starting server error", err)
	}
}
