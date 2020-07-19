package natsbroker

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
	"github.com/whale-team/whaleEcho/pkg/natspool"
)

// New construct a nats broker
func New(client natspool.Client) msgbroker.MsgBroker {
	return &NatsBroker{
		client: client,
	}
}

// NatsBroker represent message broker by using nats
type NatsBroker struct {
	client natspool.Client
}

// BindChannelWithSubject bind channel with sepcific subject
func (broker NatsBroker) BindChannelWithSubject(ctx context.Context, subject string, ch chan *nats.Msg) (entity.Subscriber, error) {
	conn, err := broker.client.SubConn()
	if err != nil {
		return nil, err
	}

	sub, err := conn.ChanSubscribe(subject, ch)
	if err != nil {
		return nil, err
	}

	if err := conn.Recycle(); err != nil {
		return nil, err
	}

	return sub, nil
}

// PublishMessage publish message to sepcific subject
func (broker NatsBroker) PublishMessage(ctx context.Context, subject string, message []byte) error {
	conn, err := broker.client.PubConn()
	if err != nil {
		return err
	}
	return conn.Publish(subject, message)
}
