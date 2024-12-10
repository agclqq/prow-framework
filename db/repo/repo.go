package repo

import "gorm.io/gorm"

func BuildMapWhere(tx *gorm.DB, data map[string]any) *gorm.DB {
	for k, v := range data {
		tx = tx.Where(k, v)
	}
	return tx
}

func ParseWhere(tx *gorm.DB, where any) *gorm.DB {
	if d, ok := where.(map[string]any); ok && d != nil {
		tx = BuildMapWhere(tx, d)
		return tx
	}
	if d, ok := where.(string); ok && d != "" {
		tx = tx.Where(d)
	}
	if where != nil {
		tx = tx.Where(where)
	}
	return tx
}

func BuildMapHaving(tx *gorm.DB, data map[string]any) *gorm.DB {
	for k, v := range data {
		tx = tx.Having(k, v)
	}
	return tx
}

func ParseHaving(tx *gorm.DB, having any) *gorm.DB {
	if d, ok := having.(map[string]any); ok && d != nil {
		tx = BuildMapHaving(tx, d)
		return tx
	}
	if d, ok := having.(string); ok && d != "" {
		tx = tx.Having(d)
	}
	if having != nil {
		tx = tx.Having(having)
	}
	return tx
}

func SelectOne(tx *gorm.DB, columns string, where any, group string, having any, order string) *gorm.DB {
	if columns != "" {
		tx = tx.Select(columns)
	}
	tx = ParseWhere(tx, where)
	if group != "" {
		tx.Group(group)
	}
	tx = ParseHaving(tx, having)
	if order != "" {
		tx.Order(order)
	}
	tx.Limit(1)
	return tx
}
func Select(tx *gorm.DB, columns string, where any, group string, having any, order string, page, pageSize int) *gorm.DB {
	if columns != "" {
		tx = tx.Select(columns)
	}
	tx = ParseWhere(tx, where)
	if group != "" {
		tx.Group(group)
	}
	tx = ParseHaving(tx, having)
	if order != "" {
		tx.Order(order)
	}
	if page > 0 && pageSize > 0 {
		tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	return tx
}

func Pagination(tx *gorm.DB, columns string, where any, group string, having any, order string, page, pageSize int) (int64, *gorm.DB) {
	var total int64
	if columns != "" {
		tx = tx.Select(columns)
	}
	tx = ParseWhere(tx, where)
	if group != "" {
		tx.Group(group)
	}
	tx = ParseHaving(tx, having)
	tx.Count(&total)
	if order != "" {
		tx.Order(order)
	}
	if page > 0 && pageSize > 0 {
		tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	return total, tx
}
