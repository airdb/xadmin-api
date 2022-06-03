package data

import (
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
)

// This interface will represent our car db
type IssueRepo interface {
	repo.Repo[IssueEntity, idkit.Id, teamworkv1.ListIssuesRequest]
}

type issueRepoDeps struct {
	fx.In
}

type issueRepo struct {
	*repokit.Repo[IssueEntity, idkit.Id, teamworkv1.ListIssuesRequest]

	deps issueRepoDeps
}

func NewIssueRepo(deps issueRepoDeps) IssueRepo {
	repo := &issueRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[IssueEntity, idkit.Id, teamworkv1.ListIssuesRequest](repo)

	return repo
}
