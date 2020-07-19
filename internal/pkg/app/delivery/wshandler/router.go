package wshandler

import "github.com/whale-team/whaleEcho/pkg/echoproto"

// SetupRoutes method used to setup routes map
func (h *Handler) SetupRoutes() {
	routeMap := make(map[echoproto.CommandType]handleFunc)

	routeMap[echoproto.CommandType_SendMessage] = h.SendMessage
	routeMap[echoproto.CommandType_JoinRoom] = h.JoinRoom
	routeMap[echoproto.CommandType_LeaveRoom] = h.LeaveRoom
	h.routeMap = routeMap
}
