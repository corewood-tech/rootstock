package server

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
)

// BinaryOnlyInterceptor rejects requests that use JSON encoding,
// enforcing binary protobuf as the only accepted content type.
func BinaryOnlyInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ct := req.Header().Get("Content-Type")
			if strings.Contains(ct, "json") {
				return nil, connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf("JSON encoding is not supported; use binary protobuf"),
				)
			}
			return next(ctx, req)
		}
	}
}
