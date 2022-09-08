package main

import (
	"github.com/dayterr/go-diploma/cmd/gophermart/handlers"
	config2 "github.com/dayterr/go-diploma/internal/config"
	"log"
	"net/http"
)

func main() {
	config, err := config2.GetConfig()
	if err != nil {
		log.Fatal("no config, can't start the program")
	}
	ah := handlers.NewAsyncHandler(config.DatabaseURI)
	r := handlers.CreateRouterWithAsyncHandler(ah)
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal("starting server error", err)
	}
}
