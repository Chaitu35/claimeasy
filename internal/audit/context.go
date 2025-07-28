package audit

import "context"

type ctxKey string

const userIDKey ctxKey = "audit_user_id"

func GetUserid(ctx context.Context) uint{
	val:= ctx.Value(userIDKey)
	if id,ok:= val.(uint); ok {
		return id
	}
	return 0
}

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}