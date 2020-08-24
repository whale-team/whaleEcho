package middleware

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"google.golang.org/protobuf/proto"
)

func WsErrorHandle(c *wsserver.Context, err error) {
	log.Error().Stack().Err(err).Msg("handler: server handle command fail")
	causeErr := errors.Cause(err)
	wsErr, ok := causeErr.(*wserrors.WsError)
	if !ok {
		wsErr = &wserrors.WsError{
			Status:  wserrors.Internal,
			Message: err.Error(),
		}
	}

	response := &echoproto.Message{
		Status:   echoproto.Status(wsErr.Status),
		Messages: []string{wsErr.Message},
		Type:     echoproto.MessageType_Response,
	}
	respData, err := proto.Marshal(response)
	if err != nil {
		log.Error().Stack().Err(wserrors.WithStack(err)).Msgf("error handler: marshal response proto failed, response:%+v", response)
		return
	}

	if err := c.WriteBinary(respData); err != nil {
		log.Error().Stack().Err(wserrors.WithStack(err)).Msg("error handler: write response binary data failed")
	}
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

func WsRecovery(c *wsserver.Context) {
	if err := recover(); err != nil {
		WsErrorHandle(c, err.(error))
	}
}
