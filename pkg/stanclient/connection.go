package stanclient

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
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
		config.ClientID = strings.Replace(config.ClientID, ".", "", -1)
	}

	proxy, err := newProxy(config)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("stan: stan client setup on %s, clientID:%s", config.Addr, config.ClientID)
	return &Client{
		conn:     proxy,
		ClientID: config.ClientID,
	}, nil
}

func buildConn(config Config, wg *sync.WaitGroup) (*nats.Conn, error) {
	wg.Add(1)
	opts := make([]nats.Option, 0, 3)
	opts = append(opts, nats.ReconnectWait(time.Duration(config.ReconnWait)*time.Second))
	opts = append(opts, nats.MaxReconnects(int(config.ReconnWait/config.ReconnDelay)))
	opts = append(opts, nats.ClosedHandler(func(*nats.Conn) {
		wg.Done()
	}))
	opts = append(opts, nats.DrainTimeout(5*time.Second))
	return nats.Connect(config.Addr, opts...)
}

func newProxy(config Config) (*stanProxy, error) {
	wg := &sync.WaitGroup{}
	natsConn, err := buildConn(config, wg)
	if err != nil {
		return nil, err
	}
	stanConn, err := stan.Connect(config.ClusterID, config.ClientID, stan.NatsConn(natsConn))
	if err != nil {
		return nil, err
	}
	return &stanProxy{
		sc: stanConn,
		wg: wg,
	}, nil
}

type stanProxy struct {
	sc stan.Conn
	wg *sync.WaitGroup
}

func (s *stanProxy) GetConn() stan.Conn {
	return s.sc
}

func (s *stanProxy) Close() error {
	err := s.sc.NatsConn().Drain()
	s.wg.Wait()
	s.sc.Close()
	return err
}
