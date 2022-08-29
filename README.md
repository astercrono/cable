# Cable

A simple, stateless, Web Socket relay server. Clients subscribe and receive messages issued to the server via a RESTful JSON PUT message.

## Usage

### Subscribing

Open up a Web Socket connection to `/sub` to receive messages.
```shell
$ websocat ws://localhost:8080/sub
{"ConnectedMs":1661738931945,"ID":"8e796cfe-8cc9-49d3-a400-0b4e753597ce","KeepAliveMs":1661739231945}
```

### Sending 

Send PUT requests of JSON data to `/sip`. These messages will be broadcasted to all connected clients.
```shell
$ curl -X PUT localhost:8080/sip -H 'Content-Type: application/json' -d '{"hello": "world"}'
$ curl -X PUT localhost:8080/sip -H 'Content-Type: application/json' -d '{"hello": "world2"}'
$ curl -X PUT localhost:8080/sip -H 'Content-Type: application/json' -d '{"hello": "world3"}'
```

### Receiving

The messages will flow down to the clients.
```shell
$ websocat ws://localhost:8080/sub # Pre-Opened Connection
{"ConnectedMs":1661738931945,"ID":"8e796cfe-8cc9-49d3-a400-0b4e753597ce","KeepAliveMs":1661739231945}  # Pre-Opened Connection
{"Client":{"ConnectedMs":1661738931945,"ID":"8e796cfe-8cc9-49d3-a400-0b4e753597ce","KeepAliveMs":1661739231945},"Message":{"Content":{"hello":"world"},"CreatedMs":1661738935827}}
{"Client":{"ConnectedMs":1661738931945,"ID":"8e796cfe-8cc9-49d3-a400-0b4e753597ce","KeepAliveMs":1661739231945},"Message":{"Content":{"hello":"world2"},"CreatedMs":1661738939083}}
{"Client":{"ConnectedMs":1661738931945,"ID":"8e796cfe-8cc9-49d3-a400-0b4e753597ce","KeepAliveMs":1661739231945},"Message":{"Content":{"hello":"world3"},"CreatedMs":1661738941755}}
```

### Keep-Alive

The configured connection timeout is 5m. Send `refresh` back to the server to refresh the connection timeout.

## Building

```shell
$ go build -o bin/cable gitlab.com/cronolabs/cable/cmd/server
$ ls bin
cable
```

## Todo

 - Support Authentication Middleware.
 - Support a default API Key system for authentication.
 - Add support for topic registration. Only receive messages for the chosen topics.
 - Config File System
