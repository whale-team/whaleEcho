package natspool

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	pool, err := New(TestConfig)
	assert.Nil(t, err)
	defer pool.Close()

	assert.Equal(t, TestConfig.PoolSize, pool.Size())

	conn, err := pool.Get()
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	assert.Equal(t, TestConfig.PoolSize-1, pool.Size())
	conn.Recycle()
	assert.Equal(t, TestConfig.PoolSize, pool.Size())

}

func TestGetTimeout(t *testing.T) {
	config := TestConfig
	config.PoolSize = 1
	pool, err := New(config)
	assert.Nil(t, err)
	defer pool.Close()

	conn, err := pool.Get()
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	conn, err = pool.Get()
	assert.NotNil(t, err)
	assert.Nil(t, conn)
}

func TestGetFlyingConn(t *testing.T) {
	pool, err := New(TestConfig)
	assert.Nil(t, err)
	defer pool.Close()
	pool.Flying(true)

	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < int(TestConfig.PoolSize); i++ {
		go func() {
			conn, err := pool.Get()
			assert.Nil(t, err)
			assert.NotNil(t, conn)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(0), pool.Size())
	conn, err := pool.Get()
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	err = conn.Recycle()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), pool.Size())
}

func TestShutdown(t *testing.T) {
	pool, err := New(TestConfig)
	assert.Nil(t, err)
	defer pool.Close()

	go func() {
		conn, _ := pool.Get()
		time.Sleep(1000 * time.Millisecond)
		conn.Recycle()
	}()

	conn, err := pool.Get()
	assert.Nil(t, err)
	conn.Subscribe("test.test", func(msg *nats.Msg) {
		fmt.Printf("rec msg: %s, from topic: %s\n", string(msg.Data), msg.Subject)
		time.Sleep(200 * time.Millisecond)
	})
	conn.Recycle()

	conn, err = pool.Get()
	assert.Nil(t, err)
	for i := 0; i < 3; i++ {
		msg := fmt.Sprintf("hello, this is a %d-th message", i+1)
		conn.Publish("test.test", []byte(msg))
	}
	conn.Recycle()

	err = pool.Shutdown()
	assert.Nil(t, err)
}
