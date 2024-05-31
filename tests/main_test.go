package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"demonstrate_orders/db/models"

	utils "demonstrate_orders/utils"

	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
)

var sc stan.Conn

func TestMain(m *testing.M) {
	var err error
	sc, err = stan.Connect("test-cluster", "test-123", stan.NatsURL("nats://nats-streaming:4223"))
	if err != nil {
		fmt.Println("Failed to connect to NATS Streaming:", err)
		os.Exit(1)
	}
	defer sc.Close()

	os.Exit(m.Run())
}

func TestOrderFlow(t *testing.T) {
	subject := "test.subject"
	order := generateRandomOrder()
	msg, err := json.Marshal(order)
	assert.NoError(t, err)
	err = sc.Publish(subject, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Published message:", string(msg))
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://app:8080/order/%s", order.OrderUID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	log.Println("Response status: ", resp.StatusCode)

	var gotOrder models.Order
	err = json.NewDecoder(resp.Body).Decode(&gotOrder)
	assert.NoError(t, err)
	assert.Equal(t, order.OrderUID, gotOrder.OrderUID)
}

func TestOrderDBPresence(t *testing.T) {
	subject := "test.subject"
	order := generateRandomOrder()
	msg, err := json.Marshal(order)
	assert.NoError(t, err)
	err = sc.Publish(subject, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Published message:", string(msg))
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	// database test
	var db *sql.DB
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", connStr)
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	defer db.Close()

	query := `SELECT order_uid FROM orders`
	rows, _ := db.Query(query)
	defer rows.Close()

	var db_order models.Order
	for rows.Next() {
		err = rows.Scan(&db_order.OrderUID)
		assert.NoError(t, err)
	}

	assert.Equal(t, order.OrderUID, db_order.OrderUID)
}

func TestNoOrder(t *testing.T) {
	url := fmt.Sprintf("http://app:8080/order/%s", "12345")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	log.Println("Response status: ", resp.StatusCode)
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
