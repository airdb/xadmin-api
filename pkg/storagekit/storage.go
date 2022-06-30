package storagekit

import (
	"fmt"
	"log"

	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

func GetDB(cfg cfg.Config, module string) (*gorm.DB, error) {
	var conCfg ConnectionConfig
	cfgData := cfg.Get(fmt.Sprintf("services.%s.database", module))
	if !cfgData.IsSet() {
		return nil, nil
	}

	log.Println("GetDB With:", cfgData.Raw())
	err := mapstructure.Decode(cfgData.Raw(), &conCfg)
	if err != nil {
		return nil, err
	}

	db, err := manager.Open(&conCfg)
	if err == nil {
		db = db.Debug()
	}

	return db, err
}
