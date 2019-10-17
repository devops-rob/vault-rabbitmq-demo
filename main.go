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

	if token == "" {
		token = "devopsrob"
	}

	if vaultaddr == "" {
		vaultaddr = "http://127.0.0.1:8200"
	}

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

	rbmqsecret, err := c.Read("rabbitmq/creds/rabbitrole")
	if err != nil {
		fmt.Println(err)
		return
	}

	vusername := rbmqsecret.Data["username"].(string)
	vpassword := rbmqsecret.Data["password"].(string)

	rmqurl := os.Getenv("AMQP_URL")
	rmqport := os.Getenv("AMQP_PORT")

	if rmqurl == "" {
		rmqurl = "127.0.0.1"
	}

	if rmqport == "" {
		rmqport = "5672"
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s", vusername, vpassword, rmqurl, rmqport)
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
		Body: []byte("Vault is absoluetly amazing"),
	}

	err = channel.Publish("events", "random-key", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	_, err = channel.QueueDeclare("vault-demo", true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	err = channel.QueueBind("vault-demo", "#", "events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	msgs, err := channel.Consume("vault-demo", "", false, false, false, false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(false)
	}

	defer connection.Close()

}
