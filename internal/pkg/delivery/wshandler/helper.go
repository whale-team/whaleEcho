package wshandler

import (
	"context"

	"github.com/vicxu416/wsserver"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"github.com/whale-team/whaleEcho/pkg/wserrors"
	"google.golang.org/protobuf/proto"
)

// UserIDKey represent context key for storing user id
type UserIDKey struct{}

// ReplyResponse function used to reply response message to websocket client
func ReplyResponse(c *wsserver.Context, status echoproto.Status, messages ...string) error {
	resp := &echoproto.Message{
		Status:   status,
		Messages: messages,
		Type:     echoproto.MessageType_Response,
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		return wserrors.Wrapf(wserrors.ErrInternal, "handler: ReplyResponse marshal response failed, err:%+v, response:%+v", err, resp)
	}
	return c.WriteBinary(data)
}

// AttachUserID attach user id into context
func AttachUserID(ctx context.Context, userUID string) context.Context {
	ctx = context.WithValue(ctx, UserIDKey{}, userUID)
	return ctx
}

// GetUserID get user id from context
func GetUserID(ctx context.Context) string {
	val := ctx.Value(UserIDKey{})
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}
