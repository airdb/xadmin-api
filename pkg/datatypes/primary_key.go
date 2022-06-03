package datatypes

import (
	"github.com/airdb/xadmin-api/pkg/idkit"
	"gorm.io/gorm"
)

type PrimaryKey struct {
	Id idkit.Id `gorm:"type:char(26);primaryKey" json:"id"`
}

func (u *PrimaryKey) String() string {
	return u.Id.String()
}

func (u *PrimaryKey) GormCondition(tx *gorm.DB) *gorm.DB {
	return tx.Where("`id` = ?", u.Id.String())
}

func (u *PrimaryKey) BeforeCreate(tx *gorm.DB) error {
	if u.Id.IsNil() {
		u.Id = idkit.New()
	}
	return nil
}
