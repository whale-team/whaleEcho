package stanclient

import (
	"context"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Subscription interface {
	Unsubscribe() error
	Close() error
	IsValid() bool
}

// Client nats stan client adapter
type Client struct {
	conn     *stanProxy
	ClientID string
}

type msgHandler func(ctx context.Context, data []byte) error

func (h msgHandler) Handle() func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		var (
			ctx       = context.Background()
			err       error
			startTime = time.Now()
			logger    zerolog.Logger
		)

		err = h(ctx, msg.Data)
		endTime := time.Now()
		logger = log.With().Fields(
			map[string]interface{}{
				"start_time": startTime,
				"end_time":   endTime,
				"latency":    strconv.FormatInt(int64(endTime.Sub(startTime)), 10),
				"data":       string(msg.Data),
			},
		).Logger()

		if err != nil {
			logger.Error().Err(err).Stack().Msgf("stan: access log, receive message on subject(%s) failed, err:%+v", msg.Subject, err)
		} else {
			logger.Info().Msgf("stan: access log, receive message on subject(%s)", msg.Subject)
		}
	}
}

// Subscribe subscribe a handler on subject
func (c *Client) Subscribe(subject string, msgCallback func(ctx context.Context, data []byte) error) error {
	_, err := c.getConn().Subscribe(subject, msgHandler(msgCallback).Handle(), c.duration())
	return err
}

func (c *Client) SubscribeQueue(subject, group string, msgCallback func(ctx context.Context, data []byte) error) error {
	_, err := c.getConn().QueueSubscribe(subject, group, msgHandler(msgCallback).Handle(), c.duration())
	return err
}

// Publish publish data to subject
func (c *Client) Publish(ctx context.Context, subject string, data []byte) error {
	return c.getConn().Publish(subject, data)
}

func (c *Client) Close() error {
	c.getConn().Close()
	c.getConn().NatsConn().Close()
	return nil
}

func (c *Client) duration() stan.SubscriptionOption {
	return stan.DurableName(c.ClientID)
}

func (c *Client) getConn() stan.Conn {
	return c.conn.GetConn()
}
