package repokit

import (
	"context"

	"gorm.io/gorm"
)

type repokitContextKey int

const dbCtxKey repokitContextKey = 1

func CreateContextDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbCtxKey, db)
}

func FromContextDB(ctx context.Context) *gorm.DB {
	val, ok := ctx.Value(dbCtxKey).(*gorm.DB)
	if !ok {
		return nil
	}
	return val
}
