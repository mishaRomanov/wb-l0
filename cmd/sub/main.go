package main

import (
	"encoding/json"
	"log"
	"net/http"
	//
	"github.com/mishaRomanov/wb-l0/internal/broker"
	"github.com/mishaRomanov/wb-l0/internal/entities"
	"github.com/mishaRomanov/wb-l0/internal/handler"
	"github.com/mishaRomanov/wb-l0/internal/storage/cache"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// consumer main
func main() {
	log.Println("Consumer service starting...")

	//creating in-memory storage
	var cache = cache.New()

	//connecting to nats server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to	nats server: %v\n", err)
	}
	// creating a new jetstream consumer
	consumer, err := broker.CreateNewConsumer(nc)

	// consuming messages from jetstream
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		//creating order struct
		order := entities.Order{}
		//unmarsharshalling json to order struct
		err := json.Unmarshal(msg.Data(), &order)
		if err != nil {
			log.Fatal(err)
		}
		//writing the data to cache
		cache.Add(order)
		//acknowledge the message
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Error while consuming: %v\n", err)
	}
	defer cc.Stop()

	generalHandler := handler.NewHandler(cache)

	//handler that return an order with given id
	http.HandleFunc("/id/{id}", generalHandler.GetByID)

	//starting a service
	serviceError := http.ListenAndServe(":8080", nil)
	if serviceError != nil {
		log.Fatalf("%v\n", serviceError)
	}
}
