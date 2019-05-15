// [Gin实践 连载十 定制 GORM Callbacks](https://segmentfault.com/a/1190000014393602)
package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

//// from https://github.com/jinzhu/gorm/blob/master/callback_create.go
// db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
//// https://github.com/jinzhu/gorm/blob/master/callback_update.go
// db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
//// https://github.com/jinzhu/gorm/blob/master/callback_delete.go
// db.Callback().Delete().Replace("gorm:delete", deleteCallback)

// for gorm
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt int64
	UpdatedAt int64
	DeletedAt *int64
}

// UpdateTimeStampForCreateCallback will set `CreatedAt`, `UpdatedAt` when creating
func UpdateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := time.Now().Unix()

		if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}

// UpdateTimeStampForUpdateCallback will set `UpdatedAt` when updating
func UpdateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", time.Now().Unix())
	}
}

// DeleteCallback used to delete data from database or set deleted_at to current time (when using with soft delete)
func DeleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedAtField, hasDeletedAtField := scope.FieldByName("DeletedAt")

		if !scope.Search.Unscoped && hasDeletedAtField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
