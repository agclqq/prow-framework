package repo

import (
	"context"

	"gorm.io/gorm"
)

type Repo[T any] interface {
	Select(ctx context.Context, distinct, columns string, where any, group string, having any, order string, page, pageSize int) ([]*T, error)
	Pagination(ctx context.Context, distinct, columns string, where any, group string, having any, order string, page, pageSize int) (int64, []*T, error)
	SelectOne(ctx context.Context, distinct, columns string, where any, group string, having any, order string) (*T, error)
	Insert(ctx context.Context, data any) error
	Update(ctx context.Context, data any, where any) (int64, error)
	Delete(ctx context.Context, where any) (int64, error)
}
type R[T any] interface {
	Select(ctx context.Context, distinct, columns string, where any, group string, having any, order string, page, pageSize int) ([]*T, *gorm.DB)
	Pagination(ctx context.Context, distinct, columns string, where any, group string, having any, order string, page, pageSize int) (int64, []*T, *gorm.DB)
	SelectOne(ctx context.Context, distinct, columns string, where any, group string, having any, order string) (*T, *gorm.DB)
	Insert(ctx context.Context, data any) *gorm.DB
	Update(ctx context.Context, data any, where any) *gorm.DB
	Delete(ctx context.Context, where any) *gorm.DB
}
type RepoImpl[T any] struct {
	db    *gorm.DB
	table string
}

func NewRepo[T any](db *gorm.DB, table string) R[T] {
	return &RepoImpl[T]{db: db, table: table}
}
func (r *RepoImpl[T]) Select(ctx context.Context, distinct, columns string, where any, group string, having any, order string, page, pageSize int) ([]*T, *gorm.DB) {
	var data []*T
	tx := r.db.WithContext(ctx).Table(r.table)
	rs := SelectQuery(tx, distinct, columns, where, group, having, order, page, pageSize).Find(&data)
	return data, rs
}
func (r *RepoImpl[T]) Pagination(ctx context.Context, distinct, columns string, where any, group string, having any, order string, page, pageSize int) (int64, []*T, *gorm.DB) {
	var data []*T
	var total int64
	tx := r.db.WithContext(ctx).Table(r.table)
	total, tx = PaginationQuery(tx, distinct, columns, where, group, having, order, page, pageSize)
	rs := tx.Find(&data)
	return total, data, rs
}
func (r *RepoImpl[T]) SelectOne(ctx context.Context, distinct, columns string, where any, group string, having any, order string) (*T, *gorm.DB) {
	var data T
	tx := r.db.WithContext(ctx).Table(r.table)
	rs := SelectOneQuery(tx, distinct, columns, where, group, having, order).Find(&data)
	return &data, rs
}
func (r *RepoImpl[T]) Insert(ctx context.Context, data any) *gorm.DB {
	return r.db.WithContext(ctx).Table(r.table).Create(data)
}

func (r *RepoImpl[T]) Update(ctx context.Context, data any, where any) *gorm.DB {
	tx := r.db.WithContext(ctx).Table(r.table)
	tx = ParseWhere(tx, where)
	return tx.Updates(data)
}
func (r *RepoImpl[T]) Delete(ctx context.Context, where any) *gorm.DB {
	var data T
	tx := r.db.WithContext(ctx).Table(r.table)
	tx = ParseWhere(tx, where)
	return tx.Delete(&data)
}

func Select[T any](ctx context.Context, db *gorm.DB, table, distinct, columns string, where any, group string, having any, order string, page, pageSize int) ([]*T, *gorm.DB) {
	var data []*T
	tx := db.WithContext(ctx).Table(table)
	rs := SelectQuery(tx, distinct, columns, where, group, having, order, page, pageSize).Find(&data)
	return data, rs
}

func Pagination[T any](ctx context.Context, db *gorm.DB, table, distinct, columns string, where any, group string, having any, order string, page, pageSize int) (int64, []*T, *gorm.DB) {
	var data []*T
	var total int64
	tx := db.WithContext(ctx).Table(table)
	total, tx = PaginationQuery(tx, distinct, columns, where, group, having, order, page, pageSize)
	rs := tx.Find(&data)
	if rs.Error != nil {
		total = 0
	}
	return total, data, rs
}
func SelectOne[T any](ctx context.Context, db *gorm.DB, table, distinct, columns string, where any, group string, having any, order string) (*T, *gorm.DB) {
	var data T
	tx := db.WithContext(ctx).Table(table)
	rs := SelectOneQuery(tx, distinct, columns, where, group, having, order).Find(&data)
	return &data, rs
}

func Insert(ctx context.Context, db *gorm.DB, table string, data any) *gorm.DB {
	return db.WithContext(ctx).Table(table).Create(data)
}

func Update(ctx context.Context, db *gorm.DB, table string, data any, where any) *gorm.DB {
	tx := db.WithContext(ctx).Table(table)
	tx = ParseWhere(tx, where)
	return tx.Updates(data)
}

func Delete[T any](ctx context.Context, db *gorm.DB, table string, where any) *gorm.DB {
	var entity T
	tx := db.WithContext(ctx).Table(table)
	tx = ParseWhere(tx, where)
	return tx.Delete(&entity)

}

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
func Assemble(tx *gorm.DB, distinct string, columns string, where any, group string, having any, order string, page, pageSize int) *gorm.DB {
	if distinct != "" {
		tx = tx.Distinct(distinct)
	}
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
	if page > 0 {
		tx.Offset((page - 1) * pageSize)
	}
	if pageSize > 0 {
		tx.Limit(pageSize)
	}
	return tx
}

func SelectOneQuery(tx *gorm.DB, distinct string, columns string, where any, group string, having any, order string) *gorm.DB {
	Assemble(tx, distinct, columns, where, group, having, order, 0, 0).Limit(1)
	return tx
}

func SelectQuery(tx *gorm.DB, distinct string, columns string, where any, group string, having any, order string, page, pageSize int) *gorm.DB {
	return Assemble(tx, distinct, columns, where, group, having, order, page, pageSize)
}

func PaginationQuery(tx *gorm.DB, distinct string, columns string, where any, group string, having any, order string, page, pageSize int) (int64, *gorm.DB) {
	var total int64
	if distinct != "" {
		tx = tx.Distinct(distinct)
	}
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
	if page > 0 {
		tx.Offset((page - 1) * pageSize)
	}
	if pageSize > 0 {
		tx.Limit(pageSize)
	}
	return total, tx
}
