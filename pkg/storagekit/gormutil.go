package storagekit

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
)

func Changed(stmt *gorm.Statement, fields ...string) (mutates []*schema.Field) {
	modelValue := stmt.ReflectValue
	switch modelValue.Kind() {
	case reflect.Slice, reflect.Array:
		modelValue = stmt.ReflectValue.Index(stmt.CurDestIndex)
	}

	selectColumns, restricted := stmt.SelectAndOmitColumns(false, true)
	changed := func(field *schema.Field) bool {
		fieldValue, _ := field.ValueOf(stmt.Context, modelValue)
		if v, ok := selectColumns[field.DBName]; (ok && v) || (!ok && !restricted) {
			if mv, mok := stmt.Dest.(map[string]interface{}); mok {
				if fv, ok := mv[field.Name]; ok {
					return !utils.AssertEqual(fv, fieldValue)
				} else if fv, ok := mv[field.DBName]; ok {
					return !utils.AssertEqual(fv, fieldValue)
				}
			} else {
				destValue := reflect.ValueOf(stmt.Dest)
				for destValue.Kind() == reflect.Ptr {
					destValue = destValue.Elem()
				}

				changedValue, zero := field.ValueOf(stmt.Context, destValue)
				if v {
					return !utils.AssertEqual(changedValue, fieldValue)
				}
				return !zero && !utils.AssertEqual(changedValue, fieldValue)
			}
		}
		return false
	}

	if len(fields) == 0 {
		for _, field := range stmt.Schema.FieldsByDBName {
			if changed(field) {
				mutates = append(mutates, field)
			}
		}
	} else {
		for _, name := range fields {
			if field := stmt.Schema.LookUpField(name); field != nil {
				if changed(field) {
					mutates = append(mutates, field)
				}
			}
		}
	}

	return
}
