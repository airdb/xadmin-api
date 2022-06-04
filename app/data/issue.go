package data

import (
	"context"
	"errors"

	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This interface will represent our car db
type IssueRepo interface {
	repo.Repo[IssueEntity, idkit.Id, teamworkv1.ListIssuesRequest]

	FindByIds(ctx context.Context, ids []string) ([]*IssueEntity, error)
	AssignToProject(ctx context.Context, ids []idkit.Id, id idkit.Id) error
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

func (repo *issueRepo) GetDB(ctx context.Context) *gorm.DB {
	return repokit.FromContextDB(ctx)
}

func (repo *issueRepo) FindByIds(ctx context.Context, ids []string) ([]*IssueEntity, error) {
	issues := []*IssueEntity{}
	err := repo.GetDB(ctx).Model(&IssueEntity{}).
		Find(&issues, "id in ?", ids).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return issues, nil
}

func (repo *issueRepo) AssignToProject(ctx context.Context, ids []idkit.Id, id idkit.Id) error {
	err := repo.GetDB(ctx).Model(&IssueEntity{}).
		Where("id in ?", ids).
		UpdateColumn("project_id", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("assign project failed")
		}
		return err
	}

	return nil
}
