package contextutils

import "context"

// Contextual key for userID
type contextKey struct {
	name string
}

var userIDContextKey = &contextKey{"userID"}

// GetUserID retrieves the userID from the context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}
