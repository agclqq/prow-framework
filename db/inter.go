package db

import "gorm.io/gorm"

type Conner interface {
	GetConn()
}

// TableStructer 表结构
type TableStructer interface {
	GetDesc(conn *gorm.DB, tableName string) [][]string
}
