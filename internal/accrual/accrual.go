package accrual

import (
	"fmt"
	"github.com/dayterr/go-diploma/internal/storage"
	"io"
	"log"
	"net/http"
	"encoding/json"
)

func (ac AccrualClient) ManagePoints(orderNumber string) {
	log.Println("got ordernumber", orderNumber)
	url := fmt.Sprintf("%s/api/orders/%s", ac.Address, orderNumber)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("getting accrual error", err)
		return
	}
	log.Println("resp is", resp)

	if resp.Status == "200 OK" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("reading response error", err)
			return
		}
		defer resp.Body.Close()

		var order storage.OrderModel
		err = json.Unmarshal(body, &order)
		if err != nil {
			log.Println("unmarshalling response error", err)
		}
		order.Number = orderNumber
		log.Println("order is", order)
		ac.Storage.UpdateOrders(order)

		switch order.Status {
		case "PROCESSED":
			userID, err := ac.Storage.FindUser(orderNumber)
			if err != nil {
				log.Println("getting user error", err)
			}
			balance, err := ac.Storage.GetBalance(int(userID))
			if err != nil {
				log.Println("getting balance error", err)
			}
			balance += order.Accrual
			ac.Storage.UpdateBalance(balance, int(userID))
		case "REGISTERED":
			log.Println("REGISTERED", orderNumber)
			ac.OrderChannel <- orderNumber
		case "PROCESSING":
			log.Println("PROCESSING", orderNumber)
			ac.OrderChannel <- orderNumber
		}
	} else {
		log.Println("other", orderNumber)
		ac.OrderChannel <- orderNumber
	}
}

func (ac AccrualClient) ReadOrderNumber() {
	log.Println("ReadOrderNumber working")
	log.Println("hoho", ac.OrderChannel)
	for ord := range ac.OrderChannel {
		log.Println("ord is", ord)
		ac.ManagePoints(ord)
	}

}