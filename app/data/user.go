package data

import (
	"context"
	"encoding/json"

	uam1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"go.uber.org/fx"
)

// This interface will represent our car db
type UserRepo interface {
	repo.Repo[UserEntity, string, uam1.ListUsersRequest]
}

type userRepoDeps struct {
	fx.In

	Config cfg.Config
}

type userRepo struct {
	*repokit.Repo[UserEntity, string, uam1.ListUsersRequest]

	deps userRepoDeps
}

func NewUserRepo(deps userRepoDeps) UserRepo {
	repo := &userRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[UserEntity, string, uam1.ListUsersRequest](repo)

	return repo
}

func (r userRepo) List(ctx context.Context, query *uam1.ListUsersRequest) ([]*UserEntity, error) {
	org := r.deps.Config.Get("xadmin.casdoor.organizationName")
	queryMap := map[string]string{
		"owner": org.String(),
	}

	url := auth.GetUrl("get-users", queryMap)

	bytes, err := auth.DoGetBytesRaw(url)
	if err != nil {
		return nil, err
	}

	var users []*UserEntity
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r userRepo) Count(ctx context.Context, query *uam1.ListUsersRequest) (int32, int32, error) {
	total, err := auth.GetUserCount("")

	return int32(total), int32(total), err
}
