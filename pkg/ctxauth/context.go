package ctxauth

import "context"

type currentUsernameCtxKey struct{}

func WithCurrentUser(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, currentUsernameCtxKey{}, userName)
}

func GetCurrentUser(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(currentUsernameCtxKey{}).(string)
	return u, ok
}
