package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/streadway/amqp"
)

func main() {

	var token = os.Getenv("TOKEN")
	var vaultaddr = os.Getenv("VAULT_ADDR")

	//vault primer
	config := &api.Config{
		Address: vaultaddr,
	}

	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	client.SetToken(token)
	c := client.Logical()

	rbmqsecret, err := c.Read("rabbitmq/creds/emrabbitmq")
	if err != nil {
		fmt.Println(err)
		return
	}

	//make connection to RabbitMQ
	rmqurl := os.Getenv("AMQP_URL")
	rmquser := os.Getenv("AMQP_USER")
	rmqport := os.Getenv("AMQP_PORT")

	if rmqurl == "" {
		rmqurl = "localhost"
	}

	if rmquser == "" {
		rmquser = "vault"
	}
	if rmqport == "" {
		rmqport = "5672"
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s", rmquser, rbmqsecret.Data, rmqurl, rmqport)
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
		Body: []byte("Hello World"),
	}

	// We publish the message to the exahange we created earlier
	err = channel.Publish("events", "random-key", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	// We create a queue named Test
	_, err = channel.QueueDeclare("test", true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue
	err = channel.QueueBind("test", "#", "events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	// We consume data in the queue named test using the channel we created in go.
	msgs, err := channel.Consume("test", "", false, false, false, false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	// We loop through the messages in the queue and print them to the console.
	// The msgs will be a go channel, not an amqp channel
	for msg := range msgs {
		//print the message to the console
		fmt.Println("message received: " + string(msg.Body))
		// Acknowledge that we have received the message so it can be removed from the queue
		msg.Ack(false)
	}

	// We close the connection after the operation has completed.
	defer connection.Close()

}
