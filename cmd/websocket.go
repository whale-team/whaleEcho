package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/vicxu416/goinfra/zlogging"
	"github.com/whale-team/whaleEcho/configs"
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker/natsbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"github.com/whale-team/whaleEcho/pkg/natspool"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
	"go.uber.org/fx"
)

// WebSocketCmd command for running websocket server
var WebSocketCmd = &cobra.Command{
	Run: runWebsocket,
	Use: "ws",
}

func runWebsocket(cmd *cobra.Command, args []string) {
	defer cmdRecover()

	config, err := configs.InitConfiguration()
	if err != nil {
		log.Error().Err(err).Msg("main: init config failed")
		os.Exit(1)
	}

	zlogging.SetupLogger(config.Log)

	var server *wsserver.SocketServer
	var natsClient natspool.Client
	var center *roomscenter.Center

	popluater := fx.New(
		fx.Supply(config, config.Nats),
		fx.Provide(natspool.NewClient, roomscenter.New),
		fx.Provide(natsbroker.New, service.New, wshandler.New, app.SetupWsServer),
		fx.Populate(&natsClient, &center, &server),
	)

	ctx := context.Background()
	if err := popluater.Start(ctx); err != nil {
		log.Error().Err(err).Msg("main: popluation process fail")
		os.Exit(1)
	}

	popluater.Stop(ctx)

	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("main: start web socket server failed!")
			os.Exit(1)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-shutdownSignal
	log.Info().Msg("main: start shutdonw process...")

	if err := natsClient.Shutdown(); err != nil {
		log.Error().Err(err).Msg("main: shutdown nats connection pools fail")
	}
	log.Info().Msg("main: shutdown nats connection pools successfully")

	center.Shutdown()
	log.Info().Msg("main: shutdown rooms center successfully")

	if err := server.Shutdown(); err != nil {
		log.Error().Err(err).Msg("main: shutdown websocket server fail")
	}
	log.Info().Msg("main: shutdown websocket server successfully")

	os.Exit(0)
}
