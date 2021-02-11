# CTF Webhook -> RMQ

Incldes script that creates a listener POSR endpoint, in request will forward the body to RMQ 

## Min requirements

go 1.14

## How to run

$ go run server.go

## How to use

start rmq
```
./rabbitmq-server
```

start ngrok to forward request to local endpoint
```
./ngrok http YOUR_PORT
```

Send post request
```
POST localhost:80
```