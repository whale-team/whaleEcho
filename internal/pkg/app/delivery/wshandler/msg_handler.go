package wshandler

import (
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/mapper"
	"github.com/whale-team/whaleEcho/pkg/bytescronv"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
)

func (h Handler) SendMessage(c *wsserver.Context, payload []byte) error {
	msg := &entity.Message{}

	if err := mapper.UnmarshalMessage(payload, msg); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "handler: SendMessage unmarshal message failed. err:%+v, payload:%s",
			err, bytescronv.BytesToString(payload))
	}

	if err := h.svc.PublishText(c.Context, msg); err != nil {
		return err
	}
	return ReplyResponse(c, echoproto.Status_OK)
}

func (h Handler) JoinRoom(c *wsserver.Context, payload []byte) error {
	room := &entity.Room{}
	user := &entity.User{}

	if err := mapper.UnmarshalRoomAndUser(payload, room, user); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "handler: SendMessage unmarshal room and user failed. err:%+v, payload:%s",
			err, bytescronv.BytesToString(payload))
	}
	user.BindConn(c.Conn)
	if err := h.svc.JoinRoom(c.Context, room.UID, user); err != nil {
		return err
	}

	return ReplyResponse(c, echoproto.Status_OK)
}

func (h Handler) LeaveRoom(c *wsserver.Context, payload []byte) error {
	room := &entity.Room{}
	user := &entity.User{}

	if err := mapper.UnmarshalRoomAndUser(payload, room, user); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "handler: SendMessage unmarshal room and user failed. err:%+v, payload:%s",
			err, bytescronv.BytesToString(payload))
	}

	if room.UID == "" {
		return wserrors.WithMessagef(wserrors.ErrInputInvalid, "room uid schould not be empty")
	}
	if user.ID == 0 {
		return wserrors.WithMessagef(wserrors.ErrInputInvalid, "user id schould not be zero")
	}

	if err := h.svc.LeaveRoom(c.Context, room.UID, user); err != nil {
		return err
	}
	return ReplyResponse(c, echoproto.Status_OK)
}
