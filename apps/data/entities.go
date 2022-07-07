package data

import (
	"time"

	dts "github.com/airdb/xadmin-api/pkg/datatypes"
	"github.com/airdb/xadmin-api/pkg/idkit"
	"github.com/airdb/xadmin-api/pkg/storagekit"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func MigratorFxOption() fx.Option {
	return fx.Options(
		fx.Supply(fx.Annotated{
			Group: storagekit.GroupMigrators,
			Target: storagekit.NewMigrator("bchm",
				&FileEntity{}, &CategoryEntity{}, &LostEntity{}, &LostStatEntity{},
			),
		}),
		fx.Supply(fx.Annotated{
			Group: storagekit.GroupMigrators,
			Target: storagekit.NewMigrator("teamwork",
				&ProjectEntity{}, &IssueEntity{},
			),
		}),
	)
}

// UserEntity is our internal representation of the car
type UserEntity struct {
	Owner     string   `xorm:"varchar(100) notnull pk" json:"owner"`
	Username  string   `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedAt dts.Time `xorm:"varchar(100)" json:"createdTime"`
	UpdatedAt dts.Time `xorm:"varchar(100)" json:"updatedTime"`

	Id                string   `xorm:"varchar(100)" json:"id"`
	Type              string   `xorm:"varchar(100)" json:"type"`
	Password          string   `xorm:"varchar(100)" json:"password"`
	PasswordSalt      string   `xorm:"varchar(100)" json:"passwordSalt"`
	DisplayName       string   `xorm:"varchar(100)" json:"displayName"`
	Avatar            string   `xorm:"varchar(255)" json:"avatar"`
	PermanentAvatar   string   `xorm:"varchar(255)" json:"permanentAvatar"`
	Email             string   `xorm:"varchar(100)" json:"email"`
	Phone             string   `xorm:"varchar(100)" json:"phone"`
	Location          string   `xorm:"varchar(100)" json:"location"`
	Address           []string `json:"address"`
	Affiliation       string   `xorm:"varchar(100)" json:"affiliation"`
	Title             string   `xorm:"varchar(100)" json:"title"`
	IdCardType        string   `xorm:"varchar(100)" json:"idCardType"`
	IdCard            string   `xorm:"varchar(100)" json:"idCard"`
	Homepage          string   `xorm:"varchar(100)" json:"homepage"`
	Bio               string   `xorm:"varchar(100)" json:"bio"`
	Tag               string   `xorm:"varchar(100)" json:"tag"`
	Region            string   `xorm:"varchar(100)" json:"region"`
	Language          string   `xorm:"varchar(100)" json:"language"`
	Gender            string   `xorm:"varchar(100)" json:"gender"`
	Birthday          string   `xorm:"varchar(100)" json:"birthday"`
	Education         string   `xorm:"varchar(100)" json:"education"`
	Score             int      `json:"score"`
	Karma             int      `json:"karma"`
	Ranking           int      `json:"ranking"`
	IsDefaultAvatar   bool     `json:"isDefaultAvatar"`
	IsOnline          bool     `json:"isOnline"`
	IsAdmin           bool     `json:"isAdmin"`
	IsGlobalAdmin     bool     `json:"isGlobalAdmin"`
	IsForbidden       bool     `json:"isForbidden"`
	IsDeleted         bool     `json:"isDeleted"`
	SignupApplication string   `xorm:"varchar(100)" json:"signupApplication"`
	Hash              string   `xorm:"varchar(100)" json:"hash"`
	PreHash           string   `xorm:"varchar(100)" json:"preHash"`

	CreatedIp    string   `xorm:"varchar(100)" json:"createdIp"`
	LastSigninAt dts.Time `xorm:"varchar(100)" json:"lastSigninTime"`
	LastSigninIp string   `xorm:"varchar(100)" json:"lastSigninIp"`

	Github   string `xorm:"varchar(100)" json:"github"`
	Google   string `xorm:"varchar(100)" json:"google"`
	QQ       string `xorm:"qq varchar(100)" json:"qq"`
	WeChat   string `xorm:"wechat varchar(100)" json:"wechat"`
	Facebook string `xorm:"facebook varchar(100)" json:"facebook"`
	DingTalk string `xorm:"dingtalk varchar(100)" json:"dingtalk"`
	Weibo    string `xorm:"weibo varchar(100)" json:"weibo"`
	Gitee    string `xorm:"gitee varchar(100)" json:"gitee"`
	LinkedIn string `xorm:"linkedin varchar(100)" json:"linkedin"`
	Wecom    string `xorm:"wecom varchar(100)" json:"wecom"`
	Lark     string `xorm:"lark varchar(100)" json:"lark"`
	Gitlab   string `xorm:"gitlab varchar(100)" json:"gitlab"`

	Ldap       string            `xorm:"ldap varchar(100)" json:"ldap"`
	Properties map[string]string `json:"properties"`
}

// PassportEntity is our internal representation of the car
type PassportEntity struct {
	Name     string
	Password string
}

// CategoryEntity is our internal representation of the car
type CategoryEntity struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	UUID string `json:"uuid"`
	// 标题
	Title string `gorm:"column:name" json:"name"`
	// 描述
	Description string `json:"description"`
}

func (e *CategoryEntity) TableName() string {
	return "tab_category"
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
	BirthedAt       time.Time `json:"birthed_at"`

	MissedCountry  string    `json:"missed_country"`
	MissedProvince string    `json:"missed_province"`
	MissedCity     string    `json:"missed_city"`
	MissedAddress  string    `json:"missed_address"`
	MissedAt       time.Time `gorm:"column:missed_at" json:"missed_at"`
	// Handler        string    `json:"handler"`
	Follower   string `json:"follower"`
	Babyid     string `json:"babyid"`
	Category   string `json:"category"`
	Height     string `json:"height"`
	SyncStatus int    `gorm:"column:syncstatus;default:0" json:"sync_status"`

	// 是否审核通过
	Audited bool `gorm:"default:false"`

	// 是否已经完结
	Done bool `gorm:"default:false"`
}

func (e *LostEntity) TableName() string {
	return "tab_lost"
}

// LostStatEntity is our internal representation of the car
type LostStatEntity struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	LostID uint
	Babyid string `json:"babyid"`

	ShareCount uint // 累计转发助力
	ShowCount  uint // 累计曝光助力
}

func (e *LostStatEntity) TableName() string {
	return "tab_lost_stat"
}

// ProjectEntity is our internal representation of the project
type ProjectEntity struct {
	dts.PrimaryKey
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Title     string `json:"title"`
	Milestone string `json:"milestone"`
	Status    string `json:"status"`
}

func (e *ProjectEntity) TableName() string {
	return "tab_projects"
}

// IssueEntity is our internal representation of the issue
type IssueEntity struct {
	dts.PrimaryKey
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	ProjectId idkit.Id `json:"project_id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
}

func (e *IssueEntity) TableName() string {
	return "tab_issues"
}
