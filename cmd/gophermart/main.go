package main

import (
	"github.com/dayterr/go-diploma/cmd/gophermart/handlers"
	"log"
	"net/http"
)

func main() {
	ah := handlers.NewAsyncHandler("")
	r := handlers.CreateRouterWithAsyncHandler(ah)
	err := http.ListenAndServe("http://localhost:8080", r)
	if err != nil {
		log.Fatal("starting server error", err)
	}
}
