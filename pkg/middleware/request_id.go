package middleware

import (
	"context"

	"github.com/vicxu416/wsserver"
)

type ReqID struct{}

func CtxWithReqID(c *wsserver.Context) context.Context {
	var (
		requsetID = c.Get(wsserver.CtxRequestID)
		key       = ReqID{}
		ctx       = c.Ctx
	)

	if requsetID != nil {
		ctx = context.WithValue(ctx, key, requsetID)
	}
	return ctx
}

func CtxGetReqID(ctx context.Context) string {
	reqID := ctx.Value(ReqID{})
	if reqID == nil {
		return ""
	}

	return reqID.(string)
}
