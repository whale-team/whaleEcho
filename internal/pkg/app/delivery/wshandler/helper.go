package wshandler

import (
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
	"google.golang.org/protobuf/proto"
)

// ReplyResponse function used to reply response message to websocket client
func ReplyResponse(c *wsserver.Context, status echoproto.Status, messages ...string) error {
	resp := &echoproto.Message{
		Status:   status,
		Messages: messages,
		Type:     echoproto.MessageType_Response,
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "handler: ReplyResponse marshal response failed, err:%+v, response:%+v", err, resp)
	}
	return c.WriteBinary(data)
}
