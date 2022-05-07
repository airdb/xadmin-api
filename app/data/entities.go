package data

import (
	"time"

	"github.com/airdb/xadmin-api/pkg/storagekit"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func MigratorFxOption() fx.Option {
	return fx.Supply(fx.Annotated{
		Group:  storagekit.GroupMigrators,
		Target: storagekit.NewMigrator("bchm", &LostEntity{}),
	})
}

// PassportEntity is our internal representation of the car
type PassportEntity struct {
	Name     string
	Password string
}

// LostEntity is our internal representation of the car
type LostEntity struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	UUID      string `json:"uuid"`
	AvatarURL string `json:"avatar_url"`
	Nickname  string `json:"nickname"`
	// 性别 0:unknown, 1:male, 2:female
	Gender uint `json:"gender"`
	// 标题
	Title string `json:"title"`
	// 子标题
	Subject    string `json:"subject"`
	Characters string `json:"characters"`
	Details    string `json:"details"`
	// 数据链接
	DataFrom        string    `json:"data_from"`
	BirthedProvince string    `json:"birthed_province"`
	BirthedCity     string    `json:"birthed_city"`
	BirthedCountry  string    `json:"birthed_country"`
	BirthedAddress  string    `json:"birthed_address"`
	BirthedAt       time.Time `gorm:"type:datetime" json:"birthed_at"`

	MissedCountry  string    `json:"missed_country"`
	MissedProvince string    `json:"missed_province"`
	MissedCity     string    `json:"missed_city"`
	MissedAddress  string    `json:"missed_address"`
	MissedAt       time.Time `gorm:"column:missed_at;type:datetime" json:"missed_at"`
	// Handler        string    `json:"handler"`
	Follower   string `json:"follower"`
	Babyid     string `json:"babyid"`
	Category   string `json:"category"`
	Height     string `json:"height"`
	SyncStatus int    `gorm:"column:syncstatus;default:0" json:"sync_status"`
}

func (e *LostEntity) TableName() string {
	return "tab_lost"
}

// func (e *LostEntity) BeforeUpdate(tx *gorm.DB) (err error) {
// 	mutates := storagekit.Changed(tx.Statement)

// 	// if no fields changed
// 	if mutates == nil {
// 		tx.Statement.BuildClauses = nil
// 	}

// 	return nil
// }
