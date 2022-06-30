package data

import (
	"context"
	"errors"

	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This interface will represent our car db
type LostStatRepo interface {
	repo.Repo[LostStatEntity, uint, bchmv1.ListLostStatsRequest]
	CreateStatByID(ctx context.Context, lostId uint) (*LostStatEntity, error)
	GetStatByID(ctx context.Context, lostId uint) (*LostStatEntity, error)
	IncreaseShare(ctx context.Context, lostId uint) error
	IncreaseShow(ctx context.Context, lostId uint) error
}

type lostStatRepoDeps struct {
	fx.In

	LostRepo LostRepo
}

type lostStatRepo struct {
	*repokit.Repo[LostStatEntity, uint, bchmv1.ListLostStatsRequest]

	deps lostStatRepoDeps
}

func NewLostStatRepo(deps lostStatRepoDeps) LostStatRepo {
	repo := &lostStatRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[LostStatEntity, uint, bchmv1.ListLostStatsRequest](repo)

	return repo
}

func (repo *lostStatRepo) GetDB(ctx context.Context) *gorm.DB {
	return repokit.FromContextDB(ctx)
}

// CreateStatByID crate a lost stat from a lost.
func (repo *lostStatRepo) CreateStatByID(ctx context.Context, lostId uint) (*LostStatEntity, error) {
	lost, err := repo.deps.LostRepo.Get(ctx, lostId)
	if err != nil {
		return nil, err
	}

	item := &LostStatEntity{
		LostID: lost.ID,
		Babyid: lost.Babyid,
	}

	err = repo.Create(ctx, item)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not exist")
		}

		return nil, errors.New("can not found record")
	}

	return item, nil
}

// GetStatByID get a lost stat by a lost id. if not exist then create it.
func (repo *lostStatRepo) GetStatByID(ctx context.Context, lostId uint) (*LostStatEntity, error) {
	item, err := repo.First(ctx, &bchmv1.ListLostStatsRequest{LostId: int32(lostId)})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repo.CreateStatByID(ctx, lostId)
		}

		return nil, errors.New("can not found record")
	}

	return item, nil
}

// IncreaseShare get a lost stat by a lost id. if not exist then create it.
func (repo *lostStatRepo) IncreaseShare(ctx context.Context, lostId uint) error {
	item := &LostStatEntity{
		LostID: lostId,
	}
	err := repo.GetDB(ctx).
		Model(item).
		Where("lost_id = ?", lostId).
		UpdateColumn("share_count", gorm.Expr("share_count + ?", 1)).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return errors.New("can not increse share")
	}

	return nil
}

// IncreaseShow get a lost stat by a lost id. if not exist then create it.
func (repo *lostStatRepo) IncreaseShow(ctx context.Context, lostId uint) error {
	item := &LostStatEntity{
		LostID: lostId,
	}
	err := repo.GetDB(ctx).
		Model(item).
		Where("lost_id = ?", lostId).
		UpdateColumn("show_count", gorm.Expr("show_count + ?", 1)).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return errors.New("can not increse show")
	}

	return nil
}
