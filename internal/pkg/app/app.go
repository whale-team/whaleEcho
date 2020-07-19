package app

import (
	"github.com/whale-team/whaleEcho/configs"
	"github.com/whale-team/whaleEcho/internal/pkg/app/delivery/wshandler"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker/natsbroker"
	"github.com/whale-team/whaleEcho/internal/pkg/app/service"
	"go.uber.org/fx"
)

// New a app container
func New(config configs.Configuration) *fx.App {
	app := fx.New(
		fx.Supply(config, config.Nats),                           // config layer
		fx.Provide(StartNatsClient, natsbroker.New),              //	 message broker layer
		fx.Provide(StartRoomsCenter, service.New, wshandler.New), // service layer
		fx.Invoke(StartWsServer),                                 // handler & server layer
	)
	return app
}
