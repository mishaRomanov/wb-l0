package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"io"
	"log"
	"os"
)

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

	// opening a models file
	file, err := os.Open("model.json")
	if err != nil {
		log.Fatal(err)
	}
	//reading file to bytes
	r, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	//publishing
	js.Publish(context.Background(), "TEST.HELLO", r)

}
