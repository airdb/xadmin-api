package data

import (
	"context"
	"log"

	commonv1 "github.com/airdb/xadmin-api/genproto/common/v1"
	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/repokit"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This interface will represent our car db
type FileRepo interface {
	repo.Repo[FileEntity, uint, commonv1.ListFilesRequest]

	GetLostByID(ctx context.Context, lostId uint) ([]*FileEntity, error)
}

type fileRepoDeps struct {
	fx.In
}

type fileRepo struct {
	*repokit.Repo[FileEntity, uint, commonv1.ListFilesRequest]

	deps fileRepoDeps
}

func NewFileRepo(deps fileRepoDeps) FileRepo {
	repo := &fileRepo{
		deps: deps,
	}
	repo.Repo = repokit.NewRepo[FileEntity, uint, commonv1.ListFilesRequest](repo)

	return repo
}

func (repo *fileRepo) GetDB(ctx context.Context) *gorm.DB {
	return repokit.FromContextDB(ctx)
}

// Create creates a new talk item.
func (repo *fileRepo) GetLostByID(ctx context.Context, lostId uint) ([]*FileEntity, error) {
	var (
		items []*FileEntity
		cnt   int64
	)

	tx := repo.GetDB(ctx).Order("sort_id desc")

	tx = tx.Where("type = 'lost' and parent_id = ?", lostId)

	d := tx.Find(&items).
		Offset(-1).
		Limit(-1).
		Count(&cnt)

	log.Println("len: ", len(items))

	return items, d.Error
}
