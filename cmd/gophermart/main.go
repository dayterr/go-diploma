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
	ah := handlers.NewAsyncHandler(config.DatabaseURI, orderChannel)
	r := handlers.CreateRouterWithAsyncHandler(ah)

	ac := accrual.NewAccrualClient(config.AccrualSystemAddress, config.DatabaseURI, orderChannel)
	go func() {
		for {
			log.Println("hey from goroutine")
			ac.ManagePoints(<-orderChannel)
		}
	}()

	err = http.ListenAndServe(config.RunAddress, r)

	if err != nil {
		log.Fatal("starting server error", err)
	}
}
