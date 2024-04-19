package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mishaRomanov/wb-l0/internal/entities"
	"github.com/mishaRomanov/wb-l0/internal/handler"
	storage "github.com/mishaRomanov/wb-l0/internal/storage/cache"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"net/http"
)

// consumer main
func main() {
	log.Println("Consumer service starting...")

	//creating in-memory storage
	var cache = storage.New()
	//connecting to nats server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to	nats server: %v\n", err)
	}
	//creating jetstream manager interface
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Error creating jetstream manager interface: %v\n", err)
	}
	//creating stream itself
	stream, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name: "TEST_STREAM",
		Subjects: []string{
			"TEST.*"}})
	if err != nil {
		log.Fatalf("Error creating stream: %v\n", err)
	}

	//creating stream consumer
	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:   "TestConsumerConsume",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("Error creating consumer: %v\n", err)
	}

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
		//print the data from stream
		fmt.Println(order)
		//acknowledge the message
		msg.Ack()
		fmt.Println("The len of cache is ", len(cache.Orders))
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
