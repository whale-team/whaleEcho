package wsserver

import (
	"github.com/rs/zerolog/log"
)

type HandleFunc func(*Context) error
type ConnBuildHandleFunc func(*Context)
type ConnCloseHandleFunc func(*Context) error
type ErrHandleFunc func(c *Context, err error)

// EchoHandle ...
func EchoHandle(c *Context) error {
	data := c.GetPayload()
	return c.WriteText(data)
}

func ErrHandle(c *Context, err error) {
	log.Error().Err(err).Msg("ws: handle unexpected error")
}

func ConnCloseHandle(c *Context) error {
	log.Info().Str("connection_id", c.ID).Msgf("ws: connection closed")
	return c.Close()
}

func ConnBuildHandle(c *Context) {
	log.Info().Str("connection_id", c.ID).Msg("ws: connection built")
}
