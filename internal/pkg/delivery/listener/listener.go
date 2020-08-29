package listener

import (
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/pkg/stanclient"
	"go.uber.org/fx"
)

// Params Listener dependencies
type Params struct {
	fx.In
	Sub *stanclient.Client
	Svc app.Servicer
}

// Listener handle publish subjects
type Listener struct {
	sub Subscriber
	svc app.Servicer
}

// New build a Listener instance
func New(params Params) Listener {
	return Listener{
		sub: params.Sub,
		svc: params.Svc,
	}
}
