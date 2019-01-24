package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var (
	// DefaultPageSize 默认页大小
	DefaultPageSize = 10
)

// LimitOffset ...
func LimitOffset(page int, perPage *int) int {
	if page <= 0 {
		page = 1
	}
	if *perPage <= 0 {
		*perPage = DefaultPageSize
	}
	return (page - 1) * *perPage
}

// SelectPageData 查询翻页数据
func SelectPageData(db *gorm.DB, outData interface{}, tableName, primaryKey string, page, perPage int, sqlCmd string, sqlValues ...interface{}) (total int64, err error) {
	expression := db.Table(tableName).Select(primaryKey).Where(sqlCmd, sqlValues...)
	if verr := expression.Count(&total).Error; gorm.IsRecordNotFoundError(verr) {
		return
	} else if verr != nil {
		err = verr
	}
	if total == 0 {
		return
	}
	offset := LimitOffset(page, &perPage)
	expression = expression.Offset(offset).Limit(perPage)
	err = db.Joins(fmt.Sprintf("inner join ? as o on o.%s = %s.%s", primaryKey, tableName, primaryKey), expression.SubQuery()).Find(outData).Error
	return
}
