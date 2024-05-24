package main

import (
	"log"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "receiver-123", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subscription, err := sc.Subscribe("test.subject", func(msg *stan.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
		log.Printf("Message Body Twice: %s %s", string(msg.Data), string(msg.Data))
		msg.Ack()
	}, stan.DurableName("my-durable"))
	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Close()
}
