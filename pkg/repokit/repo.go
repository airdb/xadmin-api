package repokit

import (
	"context"
	"errors"

	"github.com/airdb/xadmin-api/pkg/interfaces/repo"
	"github.com/airdb/xadmin-api/pkg/querykit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo[TEntity any, TKey any, TQuery any] struct {
	override interface{}
}

var _ repo.Repo[interface{}, interface{}, interface{}] = (*Repo[interface{}, interface{}, interface{}])(nil)

func NewRepo[TEntity any, TKey any, TQuery any](override interface{}) *Repo[TEntity, TKey, TQuery] {
	return &Repo[TEntity, TKey, TQuery]{override: override}
}

type GetDB interface {
	GetDB(ctx context.Context) *gorm.DB
}

func (r *Repo[TEntity, TKey, TQuery]) getDB(ctx context.Context) *gorm.DB {
	if override, ok := r.override.(GetDB); ok {
		return override.GetDB(ctx)
	}

	return FromContextDB(ctx)
}

type BuildDetailScope interface {
	BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB
}

// BuildDetailScope preload relations
func (r *Repo[TEntity, TKey, TQuery]) buildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildDetailScope); ok {
		return override.BuildDetailScope(withDetail)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

type BuildFilterScope[TQuery any] interface {
	BuildFilterScope(q *TQuery) func(db *gorm.DB) *gorm.DB
}

// BuildFilterScope filter
func (r *Repo[TEntity, TKey, TQuery]) buildFilterScope(q *TQuery) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildFilterScope[TQuery]); ok {
		return override.BuildFilterScope(q)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

type DefaultSorting interface {
	DefaultSorting() []string
}

// DefaultSorting get default sorting
func (r *Repo[TEntity, TKey, TQuery]) defaultSorting() []string {
	if override, ok := r.override.(DefaultSorting); ok {
		return override.DefaultSorting()
	}
	return nil
}

type BuildSortScope[TQuery any] interface {
	BuildSortScope(q *TQuery) func(db *gorm.DB) *gorm.DB
}

// BuildSortScope build sorting query
func (r *Repo[TEntity, TKey, TQuery]) buildSortScope(q *TQuery) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildSortScope[TQuery]); ok {
		return override.BuildSortScope(q)
	}
	f, ok := (interface{})(q).(querykit.Sort)
	if ok {
		return SortScope(f, r.defaultSorting())
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

type BuildPageScope[TQuery any] interface {
	BuildPageScope(q *TQuery) func(db *gorm.DB) *gorm.DB
}

type UpdateAssociation[TEntity any] interface {
	UpdateAssociation(ctx context.Context, entity *TEntity) error
}

// BuildPageScope page query
func (r *Repo[TEntity, TKey, TQuery]) buildPageScope(q *TQuery) func(db *gorm.DB) *gorm.DB {
	if override, ok := r.override.(BuildPageScope[TQuery]); ok {
		return override.BuildPageScope(q)
	}

	f, ok := (interface{})(q).(querykit.Page)
	if ok {
		return PageScope(f)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func (r *Repo[TEntity, TKey, TQuery]) List(ctx context.Context, query *TQuery) ([]*TEntity, error) {
	var e TEntity
	db := r.getDB(ctx).Model(&e)
	db = db.Scopes(
		r.buildFilterScope(query),
		r.buildDetailScope(false),
		r.buildSortScope(query),
		r.buildPageScope(query),
	)
	var items []*TEntity
	res := db.Find(&items)
	return items, res.Error
}

func (r *Repo[TEntity, TKey, TQuery]) First(ctx context.Context, query *TQuery) (*TEntity, error) {
	var e TEntity
	db := r.getDB(ctx).Model(&e)
	db = db.Scopes(r.buildFilterScope(query), r.buildDetailScope(true))
	var item TEntity
	err := db.First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repo[TEntity, TKey, TQuery]) Count(ctx context.Context, query *TQuery) (int32, int32, error) {
	var (
		total    int64
		filtered int64
		err      error
		e        TEntity
	)
	db := r.getDB(ctx).Model(&e)
	err = db.Count(&total).Error
	if err != nil {
		return int32(total), int32(filtered), err
	}

	db = db.Scopes(r.buildFilterScope(query))
	err = db.Count(&filtered).Error
	return int32(total), int32(filtered), err
}

func (r *Repo[TEntity, TKey, TQuery]) Get(ctx context.Context, id TKey) (*TEntity, error) {
	var entity TEntity
	err := r.getDB(ctx).Model(&entity).
		Scopes(r.buildDetailScope(true)).
		First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *Repo[TEntity, TKey, TQuery]) Create(ctx context.Context, entity *TEntity) error {
	if err := r.getDB(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repo[TEntity, TKey, TQuery]) BatchCreate(ctx context.Context, entity []*TEntity, batchSize int) error {
	if err := r.getDB(ctx).
		CreateInBatches(entity, batchSize).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repo[TEntity, TKey, TQuery]) Update(ctx context.Context, id TKey, entity *TEntity, p querykit.Fields) error {
	var e TEntity
	db := r.getDB(ctx)
	err := db.First(&e, "id = ?", id).Error
	if err != nil {
		return err
	}

	db = db.Model(&e)
	if p == nil {
		return errors.New("no effect")
	}

	pathKeep := p.Keep().Paths
	if len(pathKeep) == 0 {
		return errors.New("no effect")
	}

	db = db.Select(pathKeep)
	pathOmit := p.Omit().Paths
	if len(pathOmit) > 0 {
		db = db.Omit(pathOmit...)
	}

	if u, ok := r.override.(UpdateAssociation[TEntity]); ok {
		if err := u.UpdateAssociation(ctx, &e); err != nil {
			return err
		}
	}

	updateRet := db.Updates(entity)
	if err := updateRet.Error; err != nil {
		return err
	}
	// check row affected for concurrency
	if updateRet.RowsAffected == 0 {
		return errors.New("no effect")
	}
	return nil
}

func (r *Repo[TEntity, TKey, TQuery]) Upsert(ctx context.Context, entity *TEntity) error {
	var e TEntity
	db := r.getDB(ctx).Model(&e)
	return db.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(entity).Error
}

func (r *Repo[TEntity, TKey, TQuery]) Delete(ctx context.Context, id TKey) error {
	var entity TEntity
	err := r.getDB(ctx).Model(&entity).
		First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not found")
		}
		return err
	}
	if err := r.getDB(ctx).Delete(&entity, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func PageScope(page querykit.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == nil {
			return db
		}
		ret := db
		if page.GetPageOffset() > 0 {
			ret = db.Offset(int(page.GetPageOffset()))
		}
		if page.GetPageSize() > 0 {
			ret = db.Limit(int(page.GetPageSize()))
		}
		return ret
	}
}

// SortScope build sorting by sort and default d
func SortScope(sort querykit.Sort, d []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var s []string
		if sort != nil {
			s = sort.GetSort()
		}
		if len(s) == 0 {
			s = d
		}
		parsed := repo.ParseSort(s)
		ret := db
		if parsed != "" {
			ret = ret.Order(parsed)
		}
		return ret
	}
}
