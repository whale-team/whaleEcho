package natspool

// Client ...
type Client interface {
	PubConn() (*PoolConn, error)
	SubConn() (*PoolConn, error)
	Shutdown() error
}

func NewClient(config Config) (Client, error) {
	subPool, err := New(config)
	if err != nil {
		return nil, err
	}
	pubPool, err := New(config)
	if err != nil {
		return nil, err
	}
	pubPool.Flying(true)

	return &client{
		subPool: subPool,
		pubPool: pubPool,
	}, nil
}

type client struct {
	subPool Pool
	pubPool Pool
}

func (c *client) PubConn() (*PoolConn, error) {
	return c.pubPool.Get()
}

func (c *client) SubConn() (*PoolConn, error) {
	return c.subPool.Get()
}

func (c *client) Shutdown() error {
	var lastErr error
	if err := c.pubPool.Shutdown(); err != nil {
		lastErr = err
	}

	if err := c.subPool.Shutdown(); err != nil {
		lastErr = err
	}
	return lastErr
}
