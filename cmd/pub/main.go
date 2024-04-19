package main

import (
	"context"
	"encoding/json"
	"github.com/mishaRomanov/wb-l0/internal/entities"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"io"
	"log"
	"os"
)

// title is pretty self-explanatory
func OpenJsonFile(name string) ([]entities.Order, error) {
	//creating orders array
	var orders []entities.Order
	// opening a models file
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	//reading file to bytes
	r, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &orders)
	return orders, nil
}

// publisher main
func main() {
	//connecting to nats server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to nats server: %v\n", err)
	}

	//creating stream interface manager
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Error creating jetstream manager interface: %v\n", err)
	}

	//connecting to stream
	_, er := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name: "TEST_STREAM",
		Subjects: []string{
			"TEST.*"},
	})
	if er != nil {
		log.Fatal(er)
	}

	//opening a file
	orders, err := OpenJsonFile("model.json")
	if err != nil {
		log.Fatal(err)
	}
	//publishing
	for n := range orders {
		d, _ := json.Marshal(orders[n])
		js.Publish(context.Background(), "TEST.HELLO", d)
		log.Println("Published a message number ", n)
	}
}
