package wshandler

import (
	"github.com/rs/zerolog/log"
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/converter"
	"github.com/whale-team/whaleEcho/pkg/bytescronv"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
)

// SendMessage handlefunc used to send message to room
func (h Handler) SendMessage(c *wsserver.Context, payload []byte) error {
	msg := &entity.Message{}

	if err := converter.UnmarshalMessage(payload, msg); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "handler: SendMessage unmarshal message failed. err:%+v, payload:%s",
			err, bytescronv.BytesToString(payload))
	}

	if err := h.svc.PublishMessage(c.Ctx, msg); err != nil {
		return err
	}

	log.Debug().Msgf("handler: user(%s) publish message to room(%s)", msg.SenderName, msg.RoomUID)
	return ReplyResponse(c, echoproto.Status_OK)
}

// JoinRoom join room
func (h Handler) JoinRoom(c *wsserver.Context, payload []byte) error {
	user := &entity.User{}

	if err := converter.UnmarshalUser(payload, user); err != nil {
		return err
	}
	user.BindConn(c)
	if err := h.svc.JoinRoom(c.Ctx, user.RoomUID, user); err != nil {
		return err
	}
	return ReplyResponse(c, echoproto.Status_OK)
}
