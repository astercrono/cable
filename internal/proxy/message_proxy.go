package proxy

import (
	cc "gitlab.com/cronolabs/cable/internal/client"
	cm "gitlab.com/cronolabs/cable/internal/message"
)

var ClientTracker *cc.ClientTracker

func InitTracker() {
	if ClientTracker == nil {
		ClientTracker = cc.NewTracker()
	}
}

func RegisterClient(client *cc.Client) *cc.Client {
	return ClientTracker.RegisterClient(client)
}

func UnregisterClient(client *cc.Client) {
	ClientTracker.RemoveClient(client)
}

func Relay(msg *cm.Message) {
	clients := ClientTracker.Clients

	for i := 0; i < len(clients); i++ {
		c := clients[i]
		c.Sink <- msg
	}
}
