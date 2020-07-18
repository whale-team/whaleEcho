package natspool

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// ConnFactory factory function for building nats connection
type ConnFactory func() (*nats.Conn, error)

// Config represent nats connection configuration
type Config struct {
	Addr           string `yaml:"addr"`
	ReconnWait     int64  `yaml:"reconn_wait"`
	ReconnDelay    int64  `yaml:"reconn_delay"`
	PoolSize       int64  `yaml:"pool_size"`
	GetConnTimeout int64  `yaml:"conn_timeout"`
}

// TestConfig config for testing nats connection
var TestConfig = Config{
	Addr:           "demo.nats.io:4222",
	ReconnDelay:    1,
	ReconnWait:     int64(10 * time.Minute),
	PoolSize:       5,
	GetConnTimeout: 1,
}

func connFactory(config Config) ConnFactory {
	return func() (*nats.Conn, error) {
		return buildConn(config)
	}
}

func buildConn(config Config) (*nats.Conn, error) {
	opts := make([]nats.Option, 0, 3)
	opts = append(opts, nats.ReconnectWait(time.Duration(config.ReconnWait)*time.Second))
	opts = append(opts, nats.MaxReconnects(int(config.ReconnWait/config.ReconnDelay)))
	return nats.Connect(config.Addr, opts...)
}

// PoolConn ...
type PoolConn struct {
	*nats.Conn
	mu     sync.RWMutex
	pool   *pool
	flying bool
}

// Recycle ...
func (p *PoolConn) Recycle() error {
	p.Flush()
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.flying {
		if p.Conn != nil && !p.Conn.IsClosed() {
			if err := p.Conn.Drain(); err != nil {
				return err
			}
		}
		return nil
	}
	return p.pool.put(p.Conn)
}
