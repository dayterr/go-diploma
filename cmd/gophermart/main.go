package main

import (
	"github.com/dayterr/go-diploma/cmd/gophermart/handlers"
	"github.com/dayterr/go-diploma/internal/accrual"
	config2 "github.com/dayterr/go-diploma/internal/config"
	"github.com/dayterr/go-diploma/internal/storage"
	"log"
	"net/http"
	"syscall"
)

func main() {
	config, err := config2.GetConfig()
	if err != nil {
		log.Fatal("no config, can't start the program")
	}

	orderChannel := make(chan string)
	ah := handlers.NewAsyncHandler()
	var s storage.Storager
	s, err = storage.NewDB(config.DatabaseURI)
	if err != nil {
		log.Println("setting database error", err)
	}
	ah.Storage = s
	r := handlers.CreateRouterWithAsyncHandler(ah)

	ac := accrual.NewAccrualClient(config.AccrualSystemAddress, config.DatabaseURI, orderChannel)
	ah.AccrualClient = ac

	go ah.AccrualClient.ReadOrderNumber()

	err = http.ListenAndServe(config.RunAddress, r)

	if err != nil {
		log.Fatal("starting server error", err)
	}
}
