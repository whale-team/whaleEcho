package msgbroker

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
)

// MsgBroker provide msg broker interface for pub sub
type MsgBroker interface {
	BindChannelWithSubject(ctx context.Context, subject string, ch chan *nats.Msg) (entity.Subscriber, error)
	PublishMessage(ctx context.Context, subject string, message []byte) error
}
