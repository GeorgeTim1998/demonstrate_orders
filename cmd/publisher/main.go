package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4223")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subject := "test.subject"
	msg := []byte("MyMessage is the best!")
	err = nc.Publish(subject, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Published message:", string(msg))

	time.Sleep(time.Second)
}
