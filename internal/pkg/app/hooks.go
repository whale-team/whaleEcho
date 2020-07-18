package app

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/configs"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/roomscenter"
	"github.com/whale-team/whaleEcho/pkg/middleware"
	"github.com/whale-team/whaleEcho/pkg/natspool"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
	"go.uber.org/fx"
)

func SetupWsServer(config configs.Configuration, handler wshandler.Handler) *wsserver.SocketServer {
	server := wsserver.New()
	server.Addr = config.WsServer.Addr
	server.Port = config.WsServer.Port

	server.ErrHandler = middleware.WsErrorHandle
	server.ConnBuildHandler = middleware.WsConnBuildHandle
	server.ConnCloseHandler = middleware.WsConnCloseHandle
	server.Handler = handler.Handle
	return server
}

func StartWsServer(config configs.Configuration, handler wshandler.Handler, lc fx.Lifecycle) (*wsserver.SocketServer, error) {
	server := SetupWsServer(config, handler)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info().Msg("app: websocket server started!")
				if err := server.Start(); err != nil {
					log.Error().Err(err).Msg("app: websocket server started failed!")
					return
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			if err := server.Shutdown(); err != nil {
				return err
			}
			log.Info().Msg("app: websocket server shutdown!")
			return nil
		},
	})

	return server, nil
}

func StartRoomsCenter(broker msgbroker.MsgBroker, lc fx.Lifecycle) (*roomscenter.Center, error) {
	center, err := roomscenter.New(broker)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			center.Shutdown()
			log.Info().Msg("app: rooms center shutdown!")
			return nil
		},
	})
	return center, nil
}

func StartNatsClient(config natspool.Config, lc fx.Lifecycle) (natspool.Client, error) {
	client, err := natspool.NewClient(config)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			if err := client.Shutdown(); err != nil {
				log.Error().Err(err).Msg("app: nats connection pools shutdown failed!")
				return err
			}
			log.Info().Msg("app: nats connection pools shutdown!")
			return nil
		},
	})
	return client, nil
}
