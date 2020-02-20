#!/bin/bash

# token to use with api requests to vault
token=devopsrob

# create the rabbitmq secrets engine mount
curl \
    --header "X-Vault-Token: $token "\
    --request POST \
    --data @rmq-payload.json \
    http://127.0.0.1:8200/v1/sys/mounts/rabbitmq

# configure connectivity to rabbitmq
curl \
    --header "X-Vault-Token: $token "\
    --request POST \
    --data @rabbitmq-payload.json \
    http://127.0.0.1:8200/v1/rabbitmq/config/connection

# Create a role with the required RabbitMQ permissions
curl \
    --header "X-Vault-Token: $token" \
    --request POST \
    --data @rmq-role.json \
    http://127.0.0.1:8200/v1/rabbitmq/roles/rabbitrole

# # just a test to create creds
# curl \
#     --header "X-Vault-Token: $token" \
#     http://127.0.0.1:8200/v1/rabbitmq/creds/rabbitrole