package middleware

import (
	"github.com/vicxu416/wsserver"

	wsMiddleware "github.com/vicxu416/wsserver/middleware"
)

type WSConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

// SetupWsServer helper
func SetupWsServer(config WSConfig) *wsserver.Server {
	opts := []wsserver.Option{
		wsserver.ConnHooks(WsConnBuildHandle, WsConnCloseHandle),
		wsserver.SetLogger(wsserver.NewDefaultLogger()),
	}

	server := wsserver.New(opts...)
	server.Addr = config.Addr
	server.Port = config.Port
	server.Use(wsMiddleware.Recover(), wsMiddleware.RequestID())
	server.MsgErrorHandleFunc(WsErrorHandle)

	return server
}
