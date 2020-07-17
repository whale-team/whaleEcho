package natspool

// Client ...
type Client interface {
	PubConn() (*PoolConn, error)
	SubConn() (*PoolConn, error)
	ShutdownPool() error
}
