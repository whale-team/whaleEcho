package wshandler

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/pkg/bytescronv"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/middleware"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

// CoomandKey represent context key for storing command struct
type CoomandKey struct{}

var commandTypeMap = map[echoproto.CommandType]string{
	echoproto.CommandType_JoinRoom:    "join_room",
	echoproto.CommandType_LeaveRoom:   "leave_room",
	echoproto.CommandType_SendMessage: "send_message",
}

// type routeMap map[echoproto.CommandType]handleFunc
type handleFunc func(c *wsserver.Context, payload []byte) error

type Params struct {
	fx.In
	Svc app.Servicer
}

func New(params Params) Handler {
	return Handler{
		svc: params.Svc,
	}
}

// Handler represent handler layer for unmarshaling protobuf, routing websocket command to servicer
type Handler struct {
	svc      app.Servicer
	routeMap map[echoproto.CommandType]handleFunc
}

// Handle method used to unmarshal protobuf, log enter, leave message, and route to handlerFunc
func (h Handler) Handle(c *wsserver.Context) error {
	command := &echoproto.Command{}

	if err := proto.Unmarshal(c.Payload(), command); err != nil {
		return wserrors.Wrapf(wserrors.ErrInputInvalid, "handler: unmarshal client command failed, err:%+v. payload:%s",
			err, bytescronv.BytesToString(c.Payload()))
	}

	ctx := middleware.CtxWithReqID(c)
	c.Ctx = context.WithValue(ctx, CoomandKey{}, command)

	startAt := time.Now()
	handleFunc := h.routeMap[command.Type]

	if handleFunc == nil {
		return wserrors.ErrCommandNotFound
	}

	err := handleFunc(c, command.Payload)
	finishAt := time.Now()
	fields := map[string]interface{}{
		"command":     commandTypeMap[command.Type],
		"started_at":  startAt.Format(time.RFC3339Nano),
		"finished_at": finishAt.Format(time.RFC3339Nano),
		"cost":        strconv.Itoa(int(finishAt.Sub(startAt).Microseconds())) + "μs",
		"request_id":  c.Get(wsserver.CtxRequestID),
	}
	if err != nil {
		log.Error().Err(err).Fields(fields).Msgf("handler: access log, process command(%s) failed", commandTypeMap[command.Type])
	} else {
		log.Info().Fields(fields).Msgf("handler: access log, process command(%s) success", commandTypeMap[command.Type])
	}
	return err
}

// GetCommand method used to get Command struct from context
func (Handler) GetCommand(ctx context.Context) *echoproto.Command {
	command := ctx.Value(CoomandKey{}).(*echoproto.Command)
	return command
}
