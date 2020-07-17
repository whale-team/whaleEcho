package entity

import (
	"encoding/json"
	"sync"

	"github.com/nats-io/nats.go"
)

// Message represent websocket message
type Message struct {
	*nats.Msg
	payload Payload
	once    sync.Once
}

// Payload ...
type Payload struct {
	UserID   string
	UserName string
	Body     string
}

// Body ...
func (m *Message) Body() string {
	return m.getPayload().Body
}

// Payload ...
func (m *Message) Payload() Payload {
	return m.getPayload()
}

func (m *Message) Data() []byte {
	return m.Msg.Data
}

func (m *Message) getPayload() Payload {
	m.once.Do(func() {
		payload := Payload{}
		json.Unmarshal(m.Data(), &payload)
		m.payload = payload
	})
	return m.payload
}
