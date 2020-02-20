package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func main() {

	queueName := "pre-vault"
	url := fmt.Sprintf("amqp://guest:guest@127.0.0.1:5672")
	connection, err := amqp.Dial(url)

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	msgs, err := channel.Consume(queueName, "", false, false, false, false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(true)
	}

	defer connection.Close()

}
