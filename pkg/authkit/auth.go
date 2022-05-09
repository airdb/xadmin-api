package authkit

import (
	"context"

	"github.com/casdoor/casdoor-go-sdk/auth"
)

type casdoorContext int

const (
	casdoorUserCtx casdoorContext = 1
)

func NewContextUser(ctx context.Context, user *auth.User) context.Context {
	return context.WithValue(ctx, casdoorUserCtx, user)
}

func FromContextUser(ctx context.Context) *auth.User {
	user, ok := ctx.Value(casdoorUserCtx).(*auth.User)
	if !ok {
		return nil
	}
	return user
}
