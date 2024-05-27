package main

import (
	"demonstrate_orders/cmd/receiver"
	"demonstrate_orders/cmd/server"
	"log"
)

func main() {
	// Запуск сервера
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Error running server: %v", err)
		}
	}()

	// Запуск получателя сообщений
	go func() {
		if err := receiver.Run(); err != nil {
			log.Fatalf("Error running receiver: %v", err)
		}
	}()

	// Ожидание завершения работы приложения
	select {}
}
