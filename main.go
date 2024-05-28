package main

import (
	"demonstrate_orders/cmd/receiver"
	"demonstrate_orders/cmd/server"
	"log"
)

func main() {
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

	select {}
}
