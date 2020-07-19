package middleware

import (
	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/pkg/wsserver"
)

func WsErrorHandle(c *wsserver.Context, err error) {
	log.Error().Stack().Err(err).Msg("handler: server handle command fail")
}

func WsConnBuildHandle(c *wsserver.Context) {
	log.Logger = log.With().Fields(map[string]interface{}{
		"connection_id": c.ID,
	}).Logger()

	log.Info().Msg("ws: connection is built")
}

func WsConnCloseHandle(c *wsserver.Context) {
	if err := c.Close(); err != nil {
		log.Error().Err(err).Msg("ws: client closed connection, server close connection failed")
		return
	}
	log.Info().Msg("ws: client closed connection, server closed connection successfully")
}
