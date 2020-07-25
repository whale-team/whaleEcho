package cmd

import (
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/whale-team/whaleEcho/configs"
)

// ClientCmd command for sending ping to websocket server
var ClientCmd = &cobra.Command{
	Run: sendPing,
	Use: "client",
}

func sendPing(cmd *cobra.Command, args []string) {

	config, err := configs.InitConfiguration()
	if err != nil {
		log.Error().Err(err).Msg("main: init config failed")
		os.Exit(1)
	}

	conn, _, err := websocket.DefaultDialer.Dial("ws://"+config.WsServer.Addr+":"+config.WsServer.Port, nil)
	defer conn.Close()

	if err != nil {
		log.Error().Err(err).Msgf("main: dail websocket server failed, server:%s", config.WsServer.Addr+":"+config.WsServer.Port)
	}

	if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		log.Error().Err(err).Msg("main: write ping message failed")
		os.Exit(1)
	}
}
