package wshandler

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"github.com/whale-team/whaleEcho/pkg/bytescronv"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
	"google.golang.org/protobuf/proto"
)

// CoomandKey represent context key for storing command struct
type CoomandKey struct{}

type routeMap map[echoproto.CommandType]handleFunc
type handleFunc func(c *wsserver.Context, payload []byte) error

var commandTypeMap = map[echoproto.CommandType]string{
	echoproto.CommandType_JoinRoom:    "join_room",
	echoproto.CommandType_LeaveRoom:   "leave_room",
	echoproto.CommandType_SendMessage: "send_message",
}

// New construct a Handler
func New(svc service.Servicer) Handler {
	handler := Handler{
		svc: svc,
	}
	handler.SetupRoutes()
	return handler
}

// Handler represent handler layer for unmarshaling protobuf, routing websocket command to servicer
type Handler struct {
	svc      service.Servicer
	routeMap routeMap
}

// Handle method used to unmarshal protobuf, log enter, leave message, and route to handlerFunc
func (h Handler) Handle(c *wsserver.Context) error {
	command := &echoproto.Command{}

	if err := proto.Unmarshal(c.GetPayload(), command); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "handler: unmarshal client command failed, err:%+v. payload:%s",
			err, bytescronv.BytesToString(c.GetPayload()))
	}

	ctx := c.Context
	c.Context = context.WithValue(ctx, CoomandKey{}, command)

	log.Logger = log.With().Fields(map[string]interface{}{
		"command": commandTypeMap[command.Type],
	}).Logger()

	log.Info().Str("started_at", time.Now().Format(time.RFC3339Nano)).Msg("access log: started")
	handleFunc := h.routeMap[command.Type]
	err := handleFunc(c, command.Payload)
	log.Info().Str("finished_at", time.Now().Format(time.RFC3339Nano)).Msg("access log: finished")
	return err
}

// GetCommand method used to get Command struct from context
func (Handler) GetCommand(ctx context.Context) *echoproto.Command {
	command := ctx.Value(CoomandKey{}).(*echoproto.Command)
	return command
}
