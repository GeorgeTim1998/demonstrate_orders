package receiver

import (
	"demonstrate_orders/config"
	DB "demonstrate_orders/db"
	"demonstrate_orders/global"
	"encoding/json"
	"log"

	"os"
	"os/signal"
	"syscall"

	"demonstrate_orders/db/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/nats-io/stan.go"
)

func Run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error downloading .env: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	sc, err := stan.Connect("test-cluster", "receiver-123", stan.NatsURL(os.Getenv("NATS_URL")))
	if err != nil {
		log.Fatalf("Receiver experienced an error connecting to nats: %v, %v", err, os.Getenv("NATS_URL"))
		return err
	}
	defer sc.Close()

	subscription, err := sc.Subscribe("test.subject", func(msg *stan.Msg) {
		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			return
		}

		if order.OrderUID == "" || order.TrackNumber == "" || order.CustomerID == "" || order.Delivery.Name == "" {
			log.Printf("Error: Order is missing required fields")
			msg.Ack()
			return
		}

		if _, err := uuid.Parse(order.OrderUID); err != nil {
			log.Printf("Error: OrderUID is not in valid UUID format")
			msg.Ack()
			return
		}

		if _, err := uuid.Parse(order.Payment.Transaction); err != nil {
			log.Printf("Error: OrderUID is not in valid UUID format")
			msg.Ack()
			return
		}

		global.Mu.Lock()
		global.Cache[order.OrderUID] = order
		global.Mu.Unlock()

		log.Printf("Received message: %s", order.OrderUID)

		err = saveOrderToDB(order)
		if err != nil {
			log.Printf("Error saving order to database: %v", err)
		} else {
			log.Printf("Saved to the DB: %s", order.OrderUID)
			msg.Ack()
		}
	}, stan.DurableName("my-durable"), stan.SetManualAckMode())
	if err != nil {
		log.Fatalf("Error in subscription: %v", err)
		return err
	}
	defer subscription.Close()

	<-interrupt

	log.Println("Received interrupt signal. Closing subscription and shutting down...")
	subscription.Close()
	log.Println("Subscription closed. Exiting...")

	return nil
}

func saveOrderToDB(order models.Order) error {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading the config: %v", err)
		return err
	}

	db, err := DB.InitDB(config)
	if err != nil {
		log.Printf("Error init DB: %v", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning the connection to the db: %v", err)
		return err
	}

	ordersStmt, err := tx.Prepare(`
	INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`)
	if err != nil {
		log.Printf("Error preparing statement for orders table: %v", err)
		return err
	}
	defer ordersStmt.Close()

	_, err = ordersStmt.Exec(order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		log.Printf("Error inserting data into orders table: %v", err)
		return err
	}

	deliveriesStmt, err := tx.Prepare(`
	INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		log.Printf("Error preparing statement for deliveries table: %v", err)
		return err
	}
	defer deliveriesStmt.Close()

	_, err = deliveriesStmt.Exec(order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		log.Printf("Error inserting data into deliveries table: %v", err)
		return err
	}

	paymentsStmt, err := tx.Prepare(`
	INSERT INTO payments (transaction, order_uid, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`)
	if err != nil {
		log.Printf("Error preparing statement for payments table: %v", err)
		return err
	}
	defer paymentsStmt.Close()

	_, err = paymentsStmt.Exec(order.Payment.Transaction, order.OrderUID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		log.Printf("Error inserting data into payments table: %v", err)
		return err
	}

	itemsStmt, err := tx.Prepare(`
	INSERT INTO items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`)
	if err != nil {
		log.Printf("Error preparing statement for items table: %v", err)
		return err
	}
	defer itemsStmt.Close()

	for _, item := range order.Items {
		_, err = itemsStmt.Exec(item.ChrtID, order.OrderUID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			log.Printf("Error inserting data into items table: %v", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}
