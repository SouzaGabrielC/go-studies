package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func OpenChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func Consume(channel *amqp.Channel, queue string, out chan<- amqp.Delivery) error {
	msgs, err := channel.Consume(
		queue,
		"go-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		out <- msg
	}

	return nil
}

func Publish(channel *amqp.Channel, body string) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(body),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := channel.PublishWithContext(
		ctx,
		"amq.direct",
		"data",
		false,
		false,
		msg,
	)

	if err != nil {
		return err
	}

	return nil
}
