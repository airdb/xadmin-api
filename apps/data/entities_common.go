package data

import (
	"time"

	"gorm.io/gorm"
)

type FileEntity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	OpenID   string `json:"open_id" gorm:"column:openid"`
	UnionID  string `json:"union_id" gorm:"column:unionid"`
	UUID     string `json:"uuid"`
	Type     string `json:"type"`
	SortID   int    `json:"sort_id"`
	ParentID uint   `json:"parent_id"`
	URL      string `json:"url"`
	Status   int    `json:"status"`
}

func (e *FileEntity) TableName() string {
	return "tab_file"
}
