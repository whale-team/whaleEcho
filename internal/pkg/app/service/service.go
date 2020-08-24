package service

import (
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/pkg/stanclient"
	"go.uber.org/fx"
)

// Params service dependency
type Params struct {
	fx.In
	Dispatcher app.Dispatcher
	StanCleint *stanclient.Client
	Repo       app.Repositorier
}

// New domain service
func New(params Params) app.Servicer {
	return &service{
		dispatcher: params.Dispatcher,
		broker:     params.StanCleint,
		repo:       params.Repo,
	}
}

type service struct {
	dispatcher app.Dispatcher
	broker     app.Broker
	repo       app.Repositorier
}
