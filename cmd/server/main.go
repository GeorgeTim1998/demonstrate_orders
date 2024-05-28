package server

import (
	"encoding/json"
	"log"
	"net/http"

	"demonstrate_orders/global"

	"github.com/gorilla/mux"
)

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["orderUID"]

	global.Mu.Lock()
	order, exists := global.Cache[orderUID]
	global.Mu.Unlock()

	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/order/{orderUID}", getOrderHandler).Methods("GET")
	r.HandleFunc("/", indexHandler).Methods("GET")

	http.Handle("/", r)
	log.Println("Server is running on port 8080")

	return http.ListenAndServe(":8080", nil)
}
