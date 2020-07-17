package wshandler

import (
	"encoding/json"
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
)

type routeMap map[CommandType]handleFunc
type handleFunc func(c *wsserver.Context, payload json.RawMessage) error

func New(svc app.Servicer) Handler {
	handler := Handler{
		svc: svc,
	}
	handler.routeMap = getRouteMap(handler)
	return handler
}

type Handler struct {
	svc      app.Servicer
	routeMap routeMap
}

func (h Handler) Routing(c *wsserver.Context) error {
	command := Command{}
	if err := c.BindJSON(&command); err != nil {
		return err
	}

	handleFunc := h.routeMap[command.Type]
	return handleFunc(c, command.Payload)
}

func getRouteMap(h Handler) routeMap {
	routeMap := make(map[CommandType]handleFunc)

	routeMap[CreateRoom] = h.CreateRoom
	routeMap[EnterRoom] = h.EnterRoom
	return routeMap
}
