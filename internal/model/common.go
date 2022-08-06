package model

import (
	"fmt"
	"time"

	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/pkg/setting"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Common struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// 创建 DB 实例
func NewDBEngin(databaseSetting *setting.DatabaseSetting) (*gorm.DB, error) {
	str := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databaseSetting.DBType,
		fmt.Sprintf(str,
			databaseSetting.UserName,
			databaseSetting.Password,
			databaseSetting.Host,
			databaseSetting.DBName,
			databaseSetting.Charset,
			databaseSetting.ParseTime,
		))
	if err != nil {
		return nil, err
	}

	// 设置日志格式
	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}

	db.SingularTable(true)
	// 设置回调函数处理公共字段
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	// 如果数据库是 mysql，db.DB() 返回 sql.DB
	// 否则返回 nil
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}

// 新增方法的回调
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	fmt.Println("updateTimeStampForCreateCallback has been exected...")

	if !scope.HasError() {
		// 没有错误
		nowTime := time.Now().Unix()
		// 判断是否包含所需的字段
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			// 判断当前字段是否为空
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime) // 通过反射对字段进行赋值
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				// 为域设置值
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 更新方法的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	fmt.Println("updateTimeStampForUpdateCallback has been exected...")

	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除行为的回调
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			// 存在 deletedOn 和 IsDel 字段，执行软删除
			now := time.Now().Unix()
			// 执行SQL语句
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v, %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			// 不存在 deletedOn 和 IsDel 字段，执行硬删除
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
