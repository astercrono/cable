package message

import (
	ct "gitlab.com/cronolabs/cable/internal/types"
	"gitlab.com/cronolabs/cable/internal/util"
)

type Message struct {
	CreatedMs int64
	Content   *ct.Json
}

func NewMessage(content *ct.Json) *Message {
	return &Message{
		Content:   content,
		CreatedMs: util.CurrentMs(),
	}
}
