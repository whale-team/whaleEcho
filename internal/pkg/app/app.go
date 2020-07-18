package app

import (
	"github.com/whale-team/whaleEcho/configs"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker/natsbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"go.uber.org/fx"
)

func New(config configs.Configuration) *fx.App {
	app := fx.New(
		fx.Supply(config, config.Nats),                           // config layer
		fx.Provide(StartNatsClient, natsbroker.New),              //	 message broker layer
		fx.Provide(StartRoomsCenter, service.New, wshandler.New), // service layer
		fx.Invoke(StartWsServer),                                 // handler & server layer
	)
	return app
}

// func Populate(config configs.Configuration, natsClient *natspool.Client, center *roomscenter.Center, server *wsserver.SocketServer) error {

// 	populate := func(client natspool.Client, center2 *roomscenter.Center, server2 *wsserver.SocketServer) {
// 		natsClient = &client
// 		center = center2
// 		server = server2
// 		log.Debug().Msgf("%+v, %+v, %+v", natsClient, center, server)
// 	}

// 	ctx := context.Background()

// 	if err := app.Start(ctx); err != nil {
// 		return err
// 	}

// 	return nil
// }
