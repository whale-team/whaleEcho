package wshandler

import (
	"errors"

	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/converter"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
)

// JoinRoom join room
func (h Handler) JoinRoom(c *wsserver.Context, payload []byte) error {
	user := &entity.User{}

	if err := converter.UnmarshalUser(payload, user); err != nil {
		return err
	}
	user.BindConn(c)
	err := h.svc.JoinRoom(c.Ctx, user.RoomUID, user)

	if err != nil && errors.Is(err, service.ErrRoomOutOfLimit) {
		return ReplyResponse(c, echoproto.Status_NotAllow, "room ("+user.RoomUID+")'s members count out of limit")
	}

	if err != nil {
		return err
	}
	return ReplyResponse(c, echoproto.Status_OK)
}
