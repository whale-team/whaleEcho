package wshandler

import (
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/converter"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
)

// JoinRoom join room
func (h Handler) JoinRoom(c *wsserver.Context, payload []byte) error {
	user := &entity.User{}

	if err := converter.UnmarshalUser(payload, user); err != nil {
		return err
	}
	user.BindConn(c)
	err := h.svc.JoinRoom(c.Ctx, user.RoomUID, user)

	if err != nil && wserrors.Is(err, wserrors.ErrRoomOutOfLimit) {
		return wserrors.ErrRoomOutOfLimit
	}

	c.Set("user_uid", user.UID)

	if err != nil {
		return err
	}
	return ReplyResponse(c, echoproto.Status_OK)
}

// LeaveRoom make user leave room
func (h Handler) LeaveRoom(c *wsserver.Context, payload []byte) error {
	var (
		user = entity.User{}
		err  error
		ctx  = c.Ctx
	)

	if err = converter.UnmarshalUser(payload, &user); err != nil {
		return err
	}

	if err = h.svc.LeaveRoom(ctx, user.RoomUID, user.UID); err != nil {
		return err
	}
	return ReplyResponse(c, echoproto.Status_OK)
}
