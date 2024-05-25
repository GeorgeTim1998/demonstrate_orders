package main

import (
	"demonstrate_orders/db/models"
	utils "demonstrate_orders/utils"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subject := "test.subject"
	order := generateRandomOrder()
	msg, _ := json.Marshal(order)
	err = sc.Publish(subject, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Published message:", string(msg))

	time.Sleep(time.Second)
}

func generateRandomOrder() models.Order {
	order := models.Order{
		OrderUID:    utils.GenerateRandomUID(),
		TrackNumber: utils.GenerateRandomString(10),
		Entry:       utils.GenerateRandomString(4),
		Delivery: models.Delivery{
			Name:    utils.GenerateRandomString(10),
			Phone:   utils.GenerateRandomPhone(),
			Zip:     utils.GenerateRandomString(6),
			City:    utils.GenerateRandomString(10),
			Address: utils.GenerateRandomString(20),
			Region:  utils.GenerateRandomString(8),
			Email:   utils.GenerateRandomEmail(),
		},
		Payment: models.Payment{
			Transaction:  utils.GenerateRandomUID(),
			Currency:     "USD",
			Provider:     utils.GenerateRandomString(6),
			Amount:       rand.Intn(10000),
			PaymentDt:    int(time.Now().Unix()),
			Bank:         utils.GenerateRandomString(6),
			DeliveryCost: rand.Intn(10000),
			GoodsTotal:   rand.Intn(10000),
			CustomFee:    rand.Intn(1000),
		},
		Items: []models.Item{
			{
				ChrtID:      rand.Intn(100000),
				TrackNumber: utils.GenerateRandomString(10),
				Price:       rand.Intn(10000),
				Rid:         utils.GenerateRandomUID(),
				Name:        utils.GenerateRandomString(8),
				Sale:        rand.Intn(50),
				Size:        utils.GenerateRandomString(2),
				TotalPrice:  rand.Intn(10000),
				NmID:        rand.Intn(100000),
				Brand:       utils.GenerateRandomString(6),
				Status:      rand.Intn(5),
			},
		},
		Locale:            "en",
		InternalSignature: utils.GenerateRandomString(10),
		CustomerID:        utils.GenerateRandomString(8),
		DeliveryService:   utils.GenerateRandomString(6),
		Shardkey:          utils.GenerateRandomString(2),
		SmID:              rand.Intn(100),
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          utils.GenerateRandomString(2),
	}

	return order
}
