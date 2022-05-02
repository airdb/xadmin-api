package common

import (
	"context"

	"github.com/casdoor/casdoor-go-sdk/auth"
)

type casdoorContext int

const (
	casdoorClaimsCtx casdoorContext = 1
)

func NewCasdoorContext(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, casdoorClaimsCtx, claims)
}

func FromCurrentCasdoorContext(ctx context.Context) *auth.Claims {
	claims, ok := ctx.Value(casdoorClaimsCtx).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}
