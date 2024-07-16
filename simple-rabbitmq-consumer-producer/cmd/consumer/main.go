package main

import (
	"events-go-expert/pkg/rabbitmq"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	channel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}

	cOut := make(chan amqp.Delivery)

	go func() {
		err := rabbitmq.Consume(channel, "minha-fila", cOut)
		if err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case msg := <-cOut:
			fmt.Printf("Message Body: %s\n", string(msg.Body))
			err := msg.Ack(true)
			if err != nil {
				fmt.Println("Error on ACK: ", err)
			}
		}
	}
}
