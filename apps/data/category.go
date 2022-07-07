package data

import (
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
)

// This interface will represent our car db
type CategoryRepo interface {
	repo.Repo[CategoryEntity, uint, bchmv1.ListCategoriesRequest]
}

type categoryRepoDeps struct {
	fx.In
}

type categoryRepo struct {
	*repokit.Repo[CategoryEntity, uint, bchmv1.ListCategoriesRequest]

	deps categoryRepoDeps
}

func NewCategoryRepo(deps categoryRepoDeps) CategoryRepo {
	repo := &categoryRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[CategoryEntity, uint, bchmv1.ListCategoriesRequest](repo)

	return repo
}
