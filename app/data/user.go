package data

import (
	uam1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
)

// This interface will represent our car db
type UserRepo interface {
	repo.Repo[UserEntity, uint, uam1.ListUsersRequest]
}

type userRepoDeps struct {
	fx.In
}

type userRepo struct {
	*repokit.Repo[UserEntity, uint, uam1.ListUsersRequest]

	deps userRepoDeps
}

func NewUserRepo(deps userRepoDeps) UserRepo {
	repo := &userRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[UserEntity, uint, uam1.ListUsersRequest](repo)

	return repo
}
