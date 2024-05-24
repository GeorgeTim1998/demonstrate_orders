package main

import (
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subject := "test.subject"
	msg := []byte("MyMessage is the best!")
	err = sc.Publish(subject, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Published message:", string(msg))

	time.Sleep(time.Second)
}
