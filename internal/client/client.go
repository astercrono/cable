package client

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	cm "gitlab.com/cronolabs/cable/internal/message"
	ct "gitlab.com/cronolabs/cable/internal/types"
	cu "gitlab.com/cronolabs/cable/internal/util"
	ce "gitlab.com/cronolabs/cable/internal/websocket/event"
)

type Client struct {
	Id          string
	ConnectedMs int64
	Connection  *websocket.Conn
	KeepAliveMs int64
	Sink        chan *cm.Message
}

type ClientTracker struct {
	Clients []*Client
	sync.Mutex
}

// Client
func NewClient(conn *websocket.Conn) *Client {
	currentMs := cu.CurrentMs()
	return &Client{
		Id:          cu.NewUuid(),
		ConnectedMs: currentMs,
		KeepAliveMs: currentMs,
		Connection:  conn,
		Sink:        make(chan *cm.Message, 1),
	}
}

func (c *Client) SendMessage(msg *cm.Message) {
	err := c.Connection.WriteJSON(ct.Json{
		"Client": c.InfoMessage(),
		"Message": ct.Json{
			"Content":   msg.Content,
			"CreatedMs": msg.CreatedMs,
		},
	})

	if err != nil {
		log.Println("Message Send error: ", err)
	}
}

func (client *Client) Close() {
	client.Connection.Close()
}

func (client *Client) RefreshKeepAlive() {
	duration := time.Now().Add(5 * time.Minute)
	client.Connection.SetReadDeadline(duration)
	client.KeepAliveMs = duration.UnixNano() / 1000000
	client.SendInfoMessage()
}

func (client *Client) InfoMessage() ct.Json {
	return ct.Json{
		"ID":          client.Id,
		"ConnectedMs": client.ConnectedMs,
		"KeepAliveMs": client.KeepAliveMs,
	}
}

func (client *Client) SendInfoMessage() {
	client.Connection.WriteJSON(client.InfoMessage())
}

func (client *Client) ReadSink() {
	for {
		newMessage := <-client.Sink
		client.SendMessage(newMessage)
	}
}

func (client *Client) WaitKeepAlive() {
	client.RefreshKeepAlive()

	for {
		_, data, err := client.Connection.ReadMessage()

		if err != nil {
			log.Println("Connection KA Error: ", err)
			break
		}

		msg := string(data[:])
		log.Println("Receive: ", msg)

		if ce.IsRefresh(msg) {
			client.RefreshKeepAlive()
		} else if ce.IsClose(msg) {
			client.Close()
			break
		} else {
			log.Println("Invalid KA Message: ", msg)
		}
	}

}

// ClientTracker
func NewTracker() *ClientTracker {
	return &ClientTracker{
		Clients: make([]*Client, 0),
	}
}

func (tracker *ClientTracker) RegisterClient(client *Client) *Client {
	tracker.Lock()
	tracker.Clients = append(tracker.Clients, client)
	tracker.Unlock()
	return client
}

func (tracker *ClientTracker) RemoveClient(client *Client) {
	tracker.Lock()
	foundIndex := -1
	clientList := tracker.Clients

	for i := 0; i < len(tracker.Clients); i++ {
		if clientList[i].Id == client.Id {
			foundIndex = i
			break
		}
	}

	if foundIndex > -1 {
		tracker.Clients = append(clientList[:foundIndex], clientList[foundIndex+1:]...)
	}
	tracker.Unlock()
}
