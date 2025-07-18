package context

import (
	"context"

	"github.com/onehappyfellow/daebak-web/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		// Nothing in the context, or we wrote an invalid value
		return nil
	}
	return user
}
