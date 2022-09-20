package accrual

import (
	"github.com/dayterr/go-diploma/internal/storage"
	"log"
)

type AccrualClient struct {
	Address string
	Storage storage.Storager
	OrderChannel chan string
}

func NewAccrualClient(address, databaseURI string, orderChannel chan int) AccrualClient {
	storage, err := storage.NewDB(databaseURI)
	if err != nil {
		log.Println("setting database error", err)
	}

	return AccrualClient{
		Address: address,
		Storage: storage,
		OrderChannel: orderChannel,
	}
}
