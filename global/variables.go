package global

import (
	"demonstrate_orders/db/models"
	"sync"
)

var (
	Cache = make(map[string]models.Order)
	Mu    sync.Mutex
)
