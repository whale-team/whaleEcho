package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/vicxu416/goinfra/zlogging"
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/configs"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/listener"
	"github.com/whale-team/whaleEcho/internal/pkg/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/dispatcher"
	"github.com/whale-team/whaleEcho/internal/pkg/repository/db"
	"github.com/whale-team/whaleEcho/pkg/middleware"
	"github.com/whale-team/whaleEcho/pkg/stanclient"
	"go.uber.org/fx"
)

// WebSocketCmd command for running websocket server
var WebSocketCmd = &cobra.Command{
	Run: runWebsocket,
	Use: "ws",
}

func runWebsocket(cmd *cobra.Command, args []string) {
	defer cmdRecover()

	var (
		ctx    = context.Background()
		rms    *dispatcher.Rooms
		stan   *stanclient.Client
		server *wsserver.Server
		err    error
	)

	config, err := configs.InitConfiguration()
	if err != nil {
		log.Error().Err(err).Msg("main: init config failed")
		os.Exit(1)
	}

	zlogging.SetupLogger(config.Log)

	app := fx.New(
		fx.Supply(config),
		fx.Provide(db.NewRedis, db.New, stanclient.New),
		fx.Provide(dispatcher.NewRooms, dispatcher.New),
		fx.Provide(service.New, listener.New, wshandler.New),
		fx.Provide(middleware.SetupWsServer),
		fx.Invoke(wshandler.SetupHandler),
		fx.Populate(&rms, &stan, &server),
	)

	if err := app.Start(ctx); err != nil {
		log.Error().Err(err).Msg("main: app start fail, stop app")
		app.Stop(ctx)
		os.Exit(1)
	}

	go server.Start()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-shutdownSignal
	log.Info().Msg("main: start shutdonw process...")
	if err := app.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("main: app stop fail")
		os.Exit(1)
	}
	if err = stan.Close(); err != nil {
		log.Error().Err(err).Msgf("main: close nats streamming client failed, err:%+v", err)
	} else {
		log.Info().Msgf("main: close nats streamming client")
	}

	rms.Clear()

	if err = server.Shutdown(5 * time.Second); err != nil {
		log.Error().Err(err).Msgf("main: close websockert server failed, err:%+v", err)
	} else {
		log.Info().Msgf("main: close websocket server")
	}

	os.Exit(0)
}
