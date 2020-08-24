package stanclient

import (
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// Config nats connection config
type Config struct {
	Addr        string `yaml:"addr"`
	ClusterID   string `yaml:"cluster_id"`
	ClientID    string `yaml:"client_id"`
	ReconnWait  int64  `yaml:"reconn_wait"`
	ReconnDelay int64  `yaml:"reconn_delay"`
}

func New(config Config) (*Client, error) {
	if config.ClientID == "" {
		config.ClientID, _ = os.Hostname()
	}

	proxy, err := newProxy(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: proxy,
	}, nil
}

func buildConn(config Config) (*nats.Conn, error) {
	opts := make([]nats.Option, 0, 3)
	opts = append(opts, nats.ReconnectWait(time.Duration(config.ReconnWait)*time.Second))
	opts = append(opts, nats.MaxReconnects(int(config.ReconnWait/config.ReconnDelay)))
	return nats.Connect(config.Addr, opts...)
}

func newProxy(config Config) (*stanProxy, error) {
	natsConn, err := buildConn(config)
	if err != nil {
		return nil, err
	}
	stanConn, err := stan.Connect(config.ClusterID, config.ClientID, stan.NatsConn(natsConn))
	if err != nil {
		return nil, err
	}
	return &stanProxy{
		sc: stanConn,
	}, nil
}

type stanProxy struct {
	sc stan.Conn
}

func (s *stanProxy) GetConn() stan.Conn {
	return s.sc
}
