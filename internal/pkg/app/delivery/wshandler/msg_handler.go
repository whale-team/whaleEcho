package wshandler

import (
	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/mapper"
	"github.com/whale-team/whaleEcho/pkg/bytescronv"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
)

// SendMessage handlefunc used to send message to room
func (h Handler) SendMessage(c *wsserver.Context, payload []byte) error {
	msg := &entity.Message{}

	if err := mapper.UnmarshalMessage(payload, msg); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "handler: SendMessage unmarshal message failed. err:%+v, payload:%s",
			err, bytescronv.BytesToString(payload))
	}

	if err := h.svc.PublishText(c.Context, msg); err != nil {
		return err
	}

	log.Debug().Msgf("handler: user(%d) publish message(%s) to room(%s)", msg.Sender.ID, msg.Text, msg.Room.UID)
	return ReplyResponse(c, echoproto.Status_OK)
}

// JoinRoom handlefunc used to bind client connection with a specific room
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

	log.Debug().Msgf("handler: user(%d) join room(%s)", user.GetID(), room.UID)

	c.Context = AttachUserID(c.Context, user.GetID())
	return ReplyResponse(c, echoproto.Status_OK)
}

// LeaveRoom handlefunc used to unbind client connection, it expect client will close connection laster
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
