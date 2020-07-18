package wshandler

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
	"google.golang.org/protobuf/proto"
)

type CoomandKey struct{}

type routeMap map[echoproto.CommandType]handleFunc
type handleFunc func(c *wsserver.Context, payload []byte) error

func New(svc service.Servicer) Handler {
	handler := Handler{
		svc: svc,
	}
	handler.SetupRoutes()
	return handler
}

type Handler struct {
	svc      service.Servicer
	routeMap routeMap
}

func (h Handler) Handle(c *wsserver.Context) error {
	command := &echoproto.Command{}

	if err := proto.Unmarshal(c.GetPayload(), command); err != nil {
		return err
	}

	ctx := c.Context
	c.Context = context.WithValue(ctx, CoomandKey{}, command)

	log.Debug().Msgf("command: %+v", command)

	handleFunc := h.routeMap[command.Type]
	return handleFunc(c, command.Payload)
}

func (Handler) GetCommand(ctx context.Context) *echoproto.Command {
	command := ctx.Value(CoomandKey{}).(*echoproto.Command)
	return command
}
