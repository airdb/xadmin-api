package storagekit

import "gorm.io/gorm"

func registerPlugins(db *gorm.DB) error {
	return db.Callback().Update().Before("gorm:update").
		Register("update_no_changed_skip", updateNoChangedSkip)
}

func updateNoChangedSkip(d *gorm.DB) {
	mutates := Changed(d.Statement)
	// if no fields changed
	if mutates == nil {
		d.Statement.BuildClauses = nil
	}
}
