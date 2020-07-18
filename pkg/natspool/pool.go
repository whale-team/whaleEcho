package natspool

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// Pool represent connection pool
type Pool interface {
	Get() (*PoolConn, error)
	Close() error
	Shutdown() error
	Size() int64
	Flying(on bool)
}

// New ...
func New(config Config) (Pool, error) {
	factory := connFactory(config)
	pool := &pool{
		conns:   make(chan *nats.Conn, config.PoolSize),
		factory: factory,
		maxConn: config.PoolSize,
		wg:      &sync.WaitGroup{},
		timeout: time.Duration(config.GetConnTimeout) * time.Second,
	}

	return pool, initPool(pool)
}

func initPool(p *pool) error {
	for i := 0; i < int(p.maxConn); i++ {
		conn, err := p.factory()
		if err != nil {
			return err
		}
		conn.SetErrorHandler(func(conn *nats.Conn, sub *nats.Subscription, err error) {
			if err != nil {
				log.Error().Err(err).Msgf("nats:Async erro occur, subscription subject:%s", sub.Subject)
			}
		})
		conn.SetClosedHandler(func(conn *nats.Conn) {
			p.wg.Done()
		})
		p.wg.Add(1)
		p.put(conn)
	}
	return nil
}

type pool struct {
	wg         *sync.WaitGroup
	mu         sync.RWMutex
	conns      chan *nats.Conn
	timeout    time.Duration
	factory    ConnFactory
	maxConn    int64
	closed     bool
	flyingMode bool
}

func (p *pool) Flying(on bool) {
	p.flyingMode = on
}

func (p *pool) Get() (*PoolConn, error) {
	if p.closed {
		return nil, ErrPoolClose
	}
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.flyingMode {
		select {
		case conn := <-p.conns:
			return p.wrapConn(conn, false), nil
		default:
			conn, err := p.factory()
			if err != nil {
				return nil, err
			}
			return p.wrapConn(conn, true), nil
		}
	}

	select {
	case conn := <-p.conns:
		return p.wrapConn(conn, false), nil
	case <-time.After(p.timeout):
		return nil, ErrConnOverLimit
	}
}

func (p *pool) wrapConn(conn *nats.Conn, flying bool) *PoolConn {
	return &PoolConn{
		Conn:   conn,
		flying: flying,
		pool:   p,
	}
}

func (p *pool) put(conn *nats.Conn) error {
	if p.closed {
		return ErrPoolClose
	}

	select {
	case p.conns <- conn:
		return nil
	default:
		if !conn.IsClosed() {
			conn.Close()
		}
		return nil
	}
}

func (p *pool) Size() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return int64(len(p.conns))
}

func (p *pool) Shutdown() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var lastErr error
	var i int64 = 0
	for conn := range p.conns {
		lastErr = conn.Drain()
		i++
		if i == p.maxConn {
			close(p.conns)
		}
	}

	p.wg.Wait()
	p.factory = nil
	p.closed = true
	p.conns = nil
	return lastErr
}

func (p *pool) Close() error {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.factory = nil
	p.closed = true
	p.mu.Unlock()

	if conns == nil {
		return nil
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
	return nil
}
