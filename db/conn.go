package db

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	file2 "github.com/agclqq/prow-framework/file"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	mylogger "github.com/agclqq/prow-framework/logger"
)

var dbMap map[string]*gorm.DB

func GetConn(connName string, conf map[string]string) *gorm.DB {
	if dbMap == nil {
		dbMap = make(map[string]*gorm.DB)
	} else if v, ok := dbMap[connName]; ok {
		return v
	}

	if conf == nil {
		panic("The database configuration item of " + connName + " cannot be found")
	}
	dsn := ""
	var db *gorm.DB
	var err error
	//fmt.Println("conf:", conf)
	switch conf["driver"] {
	case "tidb":
		//through
	case "mysql":
		if conf["charset"] == "" {
			conf["charset"] = "utf8mb4"
		}
		//dsn = conf["user"] + ":" + conf["password"] + "@tcp(" + conf["host"] + ":" + conf["port"] + ")/" + conf["db"] + "?charset=" + conf["charset"] + "&parseTime=True&loc=Local"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", conf["user"], conf["password"], conf["host"], conf["port"], conf["db"], conf["charset"])
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: setLogger(conf)})
	case "pgsql":
		//dsn = "host=" + conf["host"] + " user=" + conf["user"] + " password=" + conf["password"] + " dbname=" + conf["db"] + " port=" + conf["port"] + " sslmode=disable TimeZone=Asia/Shanghai"
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", conf["host"], conf["user"], conf["password"], conf["db"], conf["port"])
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: setLogger(conf)})
	case "clickhouse":
		// dsn = "clickhouse://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&max_execution_time=60"
		dsn = fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=200ms&max_execution_time=60", conf["user"], conf["password"], conf["host"], conf["port"], conf["db"])
		db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{Logger: setLogger(conf)})
	//case "sqlite":
	//	db, err = gorm.Open(sqlite.Open(conf["db"]), &gorm.Config{Logger: setLogger(conf)})
	default:
		panic("not supported driver type：" + conf["driver"])
	}
	if err != nil {
		panic(err)
	}
	setPool(db, conf)
	dbMap[connName] = db
	return db
}

func setLogger(conf map[string]string) logger.Interface {
	if logFile, ok := conf["log"]; ok && logFile != "" {
		file, err := file2.OpenOrCreate(logFile)
		if err != nil {
			panic(err)
		}
		return mylogger.NewSQLLogger(log.New(file, "", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		})
	}
	return nil
}

func setPool(db *gorm.DB, conf map[string]string) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	if v, err := confConvToInt(conf, "maxIdle"); err == nil {
		sqlDB.SetMaxIdleConns(v)
	}

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	if v, err := confConvToInt(conf, "maxOpen"); err == nil {
		sqlDB.SetMaxOpenConns(v)
	}

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	if v, err := confConvToInt(conf, "maxLife"); err == nil {
		sqlDB.SetConnMaxLifetime(time.Duration(v) * time.Second)
	}

	// SetConnMaxIdleTime 设置了空闲连接最大生存时间。
	if v, err := confConvToInt(conf, "maxIdleTime"); err == nil {
		sqlDB.SetConnMaxIdleTime(time.Duration(v) * time.Second)
	}
}

func confConvToInt(conf map[string]string, key string) (int, error) {
	if v, ok := conf[key]; ok {
		iv, err := strconv.Atoi(v)
		return iv, err
	}
	return 0, errors.New("")
}
