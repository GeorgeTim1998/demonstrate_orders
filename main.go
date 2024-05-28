package main

import (
	"demonstrate_orders/cmd/receiver"
	"demonstrate_orders/cmd/server"
	"demonstrate_orders/global"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	err := global.LoadCacheFromDB()
	if err != nil {
		log.Fatalf("Error loading cache from database: %v", err)
	}

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Error running server: %v", err)
		}
	}()

	go func() {
		if err := receiver.Run(); err != nil {
			log.Fatalf("Error running receiver: %v", err)
		}
	}()

	<-interrupt
	log.Println("Received interrupt signal for the application. Shutting down...")
}
