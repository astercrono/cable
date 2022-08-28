package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	cl "gitlab.com/cronolabs/cable/internal/client"
	cp "gitlab.com/cronolabs/cable/internal/proxy"
)

var upgrader = websocket.Upgrader{}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error Upgrading Connection")
	}

	client := cl.NewClient(conn)
	cp.RegisterClient(client)

	defer func() {
		client.Close()
		cp.UnregisterClient(client)
	}()

	go client.ReadSink()
	client.WaitKeepAlive()
}
