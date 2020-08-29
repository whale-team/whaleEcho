package stanclient

import (
	"context"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/pkg/bytescronv"
	"github.com/whale-team/whaleEcho/pkg/middleware"
)

const (
	reqIDLen = 32
)

var (
	reqIDFlag = bytescronv.StringToBytes("$%&*&%$%&!!!1231")
	flagLen   = len(reqIDFlag)
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

		data := msg.Data
		err = h(ctx, data)
		endTime := time.Now()
		logger = log.With().Fields(
			map[string]interface{}{
				"start_time": startTime,
				"end_time":   endTime,
				"latency":    strconv.FormatInt(int64(endTime.Sub(startTime)), 10),
				"data":       string(data),
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
	return c.conn.Close()
}

func (c *Client) duration() stan.SubscriptionOption {
	return stan.DurableName(c.ClientID)
}

func (c *Client) getConn() stan.Conn {
	return c.conn.GetConn()
}

func attachReqID(ctx context.Context, data []byte) []byte {
	reqID := middleware.CtxGetReqID(ctx)
	if reqID != "" {
		data = append(data, bytescronv.StringToBytes(reqID)...)
		data = append(data, reqIDFlag...)
	}
	return data
}

func fetchReqID(data []byte) (string, []byte) {
	var requestID string
	if bytescronv.BytesToString(data[len(data)-flagLen:len(data)]) == bytescronv.BytesToString(reqIDFlag) {
		requestID = bytescronv.BytesToString((data[len(data)-flagLen-reqIDLen : len(data)-flagLen]))
		data = data[:len(data)-flagLen-reqIDLen]
	}
	return requestID, data
}
