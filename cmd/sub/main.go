package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	//
	"github.com/mishaRomanov/wb-l0/internal/broker"
	"github.com/mishaRomanov/wb-l0/internal/entities"
	"github.com/mishaRomanov/wb-l0/internal/handler"
	"github.com/mishaRomanov/wb-l0/internal/storage/cache"
	"github.com/mishaRomanov/wb-l0/internal/storage/postgres"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// consumer main
func main() {
	log.Println("Consumer service starting...")

	//creating in-memory storage
	var inMemoryCache = cache.New()

	//connecting to nats server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to	nats server: %v\n", err)
	}
	log.Println("Connection to nats server successful")

	// creating a new jetstream consumer
	consumer, err := broker.CreateNewConsumer(nc)
	if err != nil {
		log.Fatalf("Error while creating consumer: %v\n", err)
	}
	log.Println("New consumer created successfully")

	//connecting to postgres and creating a pgx instance
	pgdb, err := postgres.CreateDB()
	if err != nil {
		log.Fatalf("Error while connecting to postgres: %v\n", err)
	}

	//checking postgres connection
	var status string
	pgdb.Db.QueryRow(context.Background(), "select 'Postgres connection established'").Scan(&status)
	log.Println(status)

	//recovering cache from postgres in goroutine
	go func() {
		status := inMemoryCache.RecoverFromPostgres(pgdb)
		if !status {
			log.Fatal("Cache recovery failed")
		}
	}()

	// consuming messages from jetstream
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		//acknowledge the message
		msg.Ack()
		//creating order struct
		order := entities.Order{}
		//unmarsharshalling json to order struct
		err := json.Unmarshal(msg.Data(), &order)
		if err != nil {
			log.Printf("Error while parsing JSON. Might be unsupported type of information : %v\n", err)
			return
		}
		//writing the data to cache
		ok := inMemoryCache.Add(order)
		//if we already have the order in our memory cache
		if !ok {
			return
		}
		//writing the data to postgres
		err = pgdb.WriteData(order)
		if err != nil {
			log.Printf("Error writing data to postgres: %v\n", err)
			return
		}
		log.Printf("New order %s added.\n", order.OrderUID)
	})
	if err != nil {
		log.Fatalf("Error while consuming: %v\n", err)
	}
	defer cc.Stop()

	generalHandler := handler.NewHandler(inMemoryCache)

	//handler that return an order with given id
	http.HandleFunc("/id/{id}", generalHandler.GetByID)

	//starting a service
	serviceError := http.ListenAndServe(":8080", nil)
	if serviceError != nil {
		log.Fatalf("%v\n", serviceError)
	}
}
