package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	sc, err := nats.Connect("nats://localhost:4223")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subscription, err := sc.Subscribe("test.subject", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
		log.Printf("Message Body Twice: %s %s", string(msg.Data), string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()
}
