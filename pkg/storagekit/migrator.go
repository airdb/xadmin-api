package storagekit

import (
	"gorm.io/gorm"
)

const GroupMigrators = "sotrageKitMigrators"

type Migrator struct {
	Module  string
	Handler func(gorm.Migrator) error
}

func NewMigrator(module string, modules ...interface{}) Migrator {
	return Migrator{
		Module: module,
		Handler: func(m gorm.Migrator) error {
			return m.AutoMigrate(modules...)
		},
	}
}
