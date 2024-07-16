package main

import (
	"events-go-expert/pkg/rabbitmq"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func main() {
	channel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}

	for {
		gUuid, err := uuid.NewUUID()
		if err != nil {
			continue
		}

		message := "Hello id: " + gUuid.String()
		err = rabbitmq.Publish(channel, message)

		fmt.Println("Message published:", message)
		time.Sleep(10 * time.Second)
	}
}
