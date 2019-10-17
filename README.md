![alt text](https://www.devopsrob.com/wp-content/uploads/2019/06/DevOps-Rob-Ltd-5-e1560366034138.png "Logo Title Text 1")

# vault-rabbitmq-demo

This is a simple Go application to demonstrate how to use Hashicorp Vault to dynamically generate short-lived credetials for RabbitMQ and use these secrets invisibly at run time for a more secure coding experience.  In order for your application to connect to RabbitMQ, you need credentials in your code to make the API call, which is the problem we are solving here with the introduction of Vault.

## Getting Started

In this demo, I use docker with a compose file to spin up the two required docker containers which are:

1. RabbitMQ
2. Vault

For the sake of simplicity of this demo, I will use the default credentials for rabbitmq and start Vault in dev mode.  I will also use http accross these services.

#### DO NOT DO THIS IN ANY PRODUCTION ENVIRONMENTS.

To start the contaiers, run the following command

```bash
docker-compose up -d
```
This will bring up the two containers.  You will be able to access them at the following URLs in your web browser:

 - Vault http://127.0.0.1:8200/ui
 - RabbitMQ http://127.0.0.1:15672

To test that the default creds work you can log into the rabbitmq instance using the following username/password pair:

username: guest
password: guest

You can also log into the vault instance using the following root token

token: devopsrob

## Setting up Vault 

Vault has a built in API server which allows us to make API calls to configure and manage the server.  This allows us to codify our configuration with tools such as Terraform.  For simplicity, I have written a short bash script which makes the required API calls to:

1. Enable the RabbitMQ secrets engine
2. Configure the above secrets engine to connect to our RabbitMQ instance
3. Create a vault role which will assign the required RabbitMQ permissions to dynamic users it creates
4. Run a test to ensure that Vault can connect to RabbitMQ and provision users

The bash script references some json files which act as payloads for the API calls.

run the following command to execute this script in a *new terminal*

```bash
chmod 777 vault-setup.sh && \
./vault-setup.sh
```

After running this script, you can navigate to the RabbitMQ URL, login using the default credentials above and select the Admin tab.  Here you will see a list of users, including a user who's username starts with "token-...".  This user is created by Vault.  If you stay on this screen, you will see that the user disappaears after a while as the TTL set in the configuration has lapsed and vaukt has now revoked the user.

## Building our application

The application connects to the vault API as the first step.  The configuration paramenters that make this happen are:
 - a vault token (devopsrob)
 - the vault address (http://127.0.0.1:8200)

These values are baked into the application as it's a demo but the application allows you to set the following environment variables to overide the defaults; however, that will not be necessary for this demo.  

TOKEN
VAULT_ADDR

It's more to show how you can use enviroment variables to keep secrets out of your code.

Using these values, the application requests and receives credentials from Vault securely and stores these in variables in the code.

The next step in the application is to take the values stored in these variables and use them to make a connection to RabbitMQ.  The address and port details for connectivity are also baked into the application code but also allows you to overide these defaults with the following two Environment variables:

AMQP_URL
AMQP_PORT

Once connected, the application does the following:

1. Creates an exchange
2. Creates a message and publishes it to the exchange
3. Creates a queue called vault-demo and binds it to the exchange
4. Consumes the data in the queue and prints the consumed massage to the console.
5. Closes the connection to the RabbitMQ instance

we will build and run our application using the go CLI

```bash
go run main.go
```

after running the program, you will see the following displayed on the console:

```bash
message received: Vault is absoluetly amazing
```

## Clean up

To clean up our resources we can just run the following command:

```bash
docker-compose stop
```