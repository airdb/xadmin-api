package storagekit

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var manager = NewManger()

var gromCfg = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "xm_",
		SingularTable: true,
	},
}

type Manager struct {
	dialector map[string]gorm.Dialector
	db        map[string]*gorm.DB
}

func NewManger() *Manager {
	return &Manager{
		dialector: make(map[string]gorm.Dialector),
		db:        make(map[string]*gorm.DB),
	}
}

func (m *Manager) Open(cc *ConnectionConfig) (*gorm.DB, error) {
	if err := m.ensureDialector(cc); err != nil {
		return nil, err
	}

	if _, ok := m.db[cc.String()]; !ok {
		db, err := gorm.Open(m.dialector[cc.String()], gromCfg)
		if err != nil {
			return nil, err
		}

		if sqlDB, err := db.DB(); err == nil {
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
		}

		err = registerPlugins(db)
		if err != nil {
			return nil, err
		}

		m.db[cc.String()] = db
	}

	return m.db[cc.String()], nil
}

func (m *Manager) ensureDialector(cc *ConnectionConfig) error {
	if _, ok := m.dialector[cc.String()]; !ok {
		switch cc.Type {
		case "sqlite":
			m.dialector[cc.String()] = sqlite.Open(cc.Dsn)
		case "mysql":
			m.dialector[cc.String()] = mysql.Open(cc.Dsn)
		case "postgres":
			m.dialector[cc.String()] = postgres.Open(cc.Dsn)
		default:
			return fmt.Errorf("connection type(%s) is not supported", cc.Type)
		}
	}

	return nil
}

type ConnectionConfig struct {
	Type string
	Dsn  string
}

func (cc ConnectionConfig) String() string {
	return fmt.Sprintf("%s::%s", cc.Type, cc.Dsn)
}
