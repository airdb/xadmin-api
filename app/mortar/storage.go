package mortar

import (
	"context"
	"strings"

	"github.com/airdb/xadmin-api/pkg/repokit"
	"github.com/airdb/xadmin-api/pkg/storagekit"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/providers/groups"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func BuildStorageFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.UnaryServerInterceptors,
			Target: StorageInterceptor,
		})
}

type storageInterceptorDeps struct {
	fx.In

	Logger    log.Logger
	Config    cfg.Config
	Migrators []storagekit.Migrator `group:"sotrageKitMigrators"`
}

func StorageInterceptor(deps storageInterceptorDeps) grpc.UnaryServerInterceptor {
	log := deps.Logger.WithField("mortar", "storage")
	for _, migrator := range deps.Migrators {
		db, err := storagekit.GetDB(deps.Config, migrator.Module)
		if err != nil {
			panic(err)
		}
		if err := migrator.Handler(db.Migrator()); err != nil {
			panic(err)
		}
	}
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		pkgStr := strings.Split(strings.Trim(info.FullMethod, "/"), "/")[0]
		pkgInfo := strings.Split(pkgStr, ".")

		db, err := storagekit.GetDB(deps.Config, pkgInfo[0])
		if err != nil {
			log.Warn(ctx, "cat get db with %s's database", pkgStr)
			return handler(ctx, req)
		}

		return handler(repokit.CreateContextDB(ctx, db), req)
	}
}
