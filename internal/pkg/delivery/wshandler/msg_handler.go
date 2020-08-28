package wshandler

import (
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
)

func (h Handler) PublishMessage(c *wsserver.Context, payload []byte) error {
	if err := h.svc.PublishRoomMessage(c.Ctx, payload); err != nil {
		return err
	}
	return ReplyResponse(c, echoproto.Status_OK)
}
