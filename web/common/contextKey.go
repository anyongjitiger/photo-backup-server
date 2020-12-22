package common

import (
	"context"

	"github.com/julienschmidt/httprouter"
)

type contextKey string

var (
    contextKeyParams = contextKey("params")
)
func WithParams(ctx context.Context, params httprouter.Params) context.Context {
	return context.WithValue(ctx, contextKeyParams, params)
}
func ParamsFromContext(ctx context.Context) (httprouter.Params, bool) {
	params, ok := ctx.Value(contextKeyParams).(httprouter.Params)
	return params, ok
}