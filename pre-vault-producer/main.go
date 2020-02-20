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

	message := amqp.Publishing{
		Body: []byte("Vault is absoluetly amazing, Why are we not using it yet?"),
	}

	err = channel.Publish("events", "random-key", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	err = channel.QueueBind(queueName, "#", "events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	defer connection.Close()

}
