# golang nats.io websocket

transports full nats protocol

## nats 

on nats server run:
- start gnatsd server then
- cd ws-to-tcp && go run main.go -token test -bind :8888

on nats client run:
- cd nats-client && go run main.go -to ws://ws.domain.tld:8888/mq?token=test -nats-user test -nats-pass test # user and pass from nats.server config (can be empty)

## nats streaming

on nats streaming server run:
- start gnatsd server then
- cd ws-to-tcp && go run main.go -token test -bind :8888

on nats client run:
- cd nats-streaming-client && go run main.go -to ws://ws.domain.tld:8888/mq?token=test -nats-user test -nats-pass test # user and pass from nats.server config (can be empty)



works behind nginx via ssl+vhost (wss)
