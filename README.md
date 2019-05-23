# golang nats.io websocket

on nats server run:
- start gnatsd server then
- cd ws-to-tcp && go run main.go -token test -bind :8888

on client run:
- cd client && go run main.go -to ws://ws.domain.tld:8888/mq?token=test -nats-user test -nats-pass test # user and pass from nats.server config (can be not set)
