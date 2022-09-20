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
	url := fmt.Sprintf("http://%s/api/orders/%d", ac.Address, orderNumber)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("getting accrual error", err)
		return
	}

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
			ac.OrderChannel <- orderNumber
		case "PROCESSING":
			ac.OrderChannel <- orderNumber
		}
	} else {
		ac.OrderChannel <- orderNumber
	}

}
