package global

import (
	"database/sql"
	"demonstrate_orders/config"
	"demonstrate_orders/db"
	"demonstrate_orders/db/models"
	"log"

	"sync"
)

var (
	Cache = make(map[string]models.Order)
	Mu    sync.Mutex
)

func LoadCacheFromDB() error {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading the config: %v", err)
		return err
	}

	database, err := db.InitDB(config)
	if err != nil {
		return err
	}
	defer database.Close()

	orders, err := fetchOrders(database)
	if err != nil {
		return err
	}

	Mu.Lock()
	for _, order := range orders {
		Cache[order.OrderUID] = order
	}
	Mu.Unlock()

	log.Println("Cache initialized from database")
	return nil
}

func fetchOrders(db *sql.DB) ([]models.Order, error) {
	query := `
	SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
		o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, 
		p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
		i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, 
		i.nm_id, i.brand, i.status
	FROM orders o
	INNER JOIN deliveries d ON o.order_uid = d.order_uid
	INNER JOIN payments p ON o.order_uid = p.order_uid
	LEFT JOIN items i ON o.order_uid = i.order_uid`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[string]*models.Order)
	for rows.Next() {
		var order models.Order
		var item models.Item

		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
			&order.Delivery.Region, &order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}

		if existingOrder, exists := orderMap[order.OrderUID]; exists {
			existingOrder.Items = append(existingOrder.Items, item)
		} else {
			order.Items = []models.Item{item}
			orderMap[order.OrderUID] = &order
		}
	}

	var orders []models.Order
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}
