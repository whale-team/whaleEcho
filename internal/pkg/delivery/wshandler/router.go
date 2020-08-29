package wshandler

import (
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
)

// SetupRoutes method used to setup routes map
func (h *Handler) SetupRoutes() {
	routeMap := make(map[echoproto.CommandType]handleFunc)

	routeMap[echoproto.CommandType_SendMessage] = h.PublishMessage
	routeMap[echoproto.CommandType_JoinRoom] = h.JoinRoom
	routeMap[echoproto.CommandType_LeaveRoom] = h.LeaveRoom
	h.routeMap = routeMap
}

// SetupHandler bind msg handler function on websocker server
func SetupHandler(serv *wsserver.Server, h Handler) error {
	h.SetupRoutes()
	serv.MsgHandlerFunc = h.Handle
	return nil
}
