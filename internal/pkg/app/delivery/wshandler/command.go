package wshandler

import (
	"encoding/json"
)

type CommandType int8

var (
	CreateRoom CommandType = 100
	EnterRoom  CommandType = 101
)

// Command ...
type Command struct {
	Type    CommandType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
