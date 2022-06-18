package data

import (
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
)

// This interface will represent our car db
type LostRepo interface {
	repo.Repo[LostEntity, uint, bchmv1.ListLostsRequest]
}

type lostRepoDeps struct {
	fx.In
}

type lostRepo struct {
	*repokit.Repo[LostEntity, uint, bchmv1.ListLostsRequest]

	deps lostRepoDeps
}

func NewLostRepo(deps lostRepoDeps) LostRepo {
	repo := &lostRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[LostEntity, uint, bchmv1.ListLostsRequest](repo)

	return repo
}
