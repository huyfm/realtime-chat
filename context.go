package rtc

import "context"

type contextKey int

const UserContextKey contextKey = 1

func UserInContext(ctx context.Context) *User {
	u, ok := ctx.Value(UserContextKey).(*User)
	if !ok {
		return nil
	}
	return u
}

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}
