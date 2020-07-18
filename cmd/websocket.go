package cmd

import (
	"github.com/whale-team/whaleEcho/pkg/wsserver"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// WebSocketCmd command for running websocket server
var WebSocketCmd = &cobra.Command{
	Run: runWebsocket,
	Use: "ws",
}

func runWebsocket(cmd *cobra.Command, args []string) {
	defer cmdRecover()

	logMiddleware := func(handler wsserver.HandleFunc) wsserver.HandleFunc {

		return func(c *wsserver.Context) error {
			log.Debug().Msg("rec message")
			return handler(c)
		}
	}

	server := wsserver.New()
	server.Handler = logMiddleware(wsserver.EchoHandle)
	server.ErrHandler = func(c *wsserver.Context, err error) {
		log.Error().Stack().Err(err).Msg("read message failed")
	}
	server.ConnBuildHandleFunc = wsserver.ConnBuildHandle
	server.ConnCloseHandler = func(c *wsserver.Context) error {
		log.Debug().Msg("conn closed")
		return nil
	}

	if err := server.ListenAndServe("", "3333"); err != nil {
		log.Error().Err(err).Msg("cmd: WebSocket cmd server startup failed")
	}
}
