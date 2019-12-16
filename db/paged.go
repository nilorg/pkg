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
func SelectPageData(db *gorm.DB, outData interface{}, primaryKey string, page, perPage int, sqlCmd string, sqlValues ...interface{}) (total int64, err error) {
	expression := db.Model(outData).Select(primaryKey).Where(sqlCmd, sqlValues...)
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
	// 阿里巴巴Java开发手册终极版v1.3.0
	// MySQL数据库->(二)索引规约->7. 【推荐】利用延迟关联或者子查询优化超多分页场景。
	//说明:MySQL 并不是跳过 offset 行，而是取 offset+N 行，然后返回放弃前 offset 行，返回 N 行，那当 offset 特别大的时候，效率就非常的低下，要么控制返回的总页数，要么对超过 特定阈值的页数进行 SQL 改写。
	//正例:先快速定位需要获取的 id 段，然后再关联:
	//	SELECT a.* FROM 表 1 a, (select id from 表 1 where 条件 LIMIT 100000,20 ) b where a.id=b.id

	err = db.Raw(fmt.Sprintf("SELECT a.* FROM `%s` a, ? b WHERE a.%s = b.%s", db.NewScope(outData).TableName(), primaryKey, primaryKey), expression.SubQuery()).Scan(outData).Error
	//err = db.Joins(fmt.Sprintf("inner join ? as o on o.%s = %s.%s", primaryKey, tableName, primaryKey), expression.SubQuery()).Find(outData).Error
	return
}
