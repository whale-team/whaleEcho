package wshandler

import (
	"encoding/json"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
)

func (h Handler) CreateRoom(c *wsserver.Context, payload json.RawMessage) error {
	return nil
}

func (h Handler) EnterRoom(c *wsserver.Context, payload json.RawMessage) error {
	return nil
}
