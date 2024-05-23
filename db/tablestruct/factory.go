package tablestruct

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/agclqq/prow-framework/prowjob/command"

	"github.com/agclqq/prow-framework/db"
)

// GetTableStruct 获取表结构
func GetTableStruct(configName string, config map[string]string) [][]string {
	conn := db.GetConn(configName, config)
	if conn == nil {
		fmt.Printf(command.NO_DB_CONFIG, configName)
		return nil
	}
	var desc [][]string
	switch conn.Name() {
	case "sqlite":
		//fallthrough
	case "mysql":
		desc = getTableStruct(Mysql{}, conn, tableName)
	}
	return desc
}
func getTableStruct(ts db.TableStructer, conn *gorm.DB, tableName string) [][]string {
	return ts.GetDesc(conn, tableName)
}
