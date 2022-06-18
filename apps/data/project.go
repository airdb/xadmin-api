package data

import (
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
)

// This interface will represent our car db
type ProjectRepo interface {
	repo.Repo[ProjectEntity, idkit.Id, teamworkv1.ListProjectsRequest]
}

type projectRepoDeps struct {
	fx.In
}

type projectRepo struct {
	*repokit.Repo[ProjectEntity, idkit.Id, teamworkv1.ListProjectsRequest]

	deps projectRepoDeps
}

func NewProjectRepo(deps projectRepoDeps) ProjectRepo {
	repo := &projectRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[ProjectEntity, idkit.Id, teamworkv1.ListProjectsRequest](repo)

	return repo
}
