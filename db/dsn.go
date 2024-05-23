package db

import (
	"errors"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cast"
)

type Config struct {
	Type      string
	Host      string
	Port      string
	User      string
	Password  string
	Dbname    string
	DbAlias   string
	Charset   string
	Collation string
	TimeZone  string
}

// Dsn 数据库连接字符串
func Dsn(config *Config) string {
	switch config.Type {
	case "mysql":
		return config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Dbname + "?charset=" + config.Charset + "&parseTime=True&loc=" + config.TimeZone
	case "pgsql":
		return "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.Dbname + " port=" + config.Port + " sslmode=disable TimeZone=" + config.TimeZone
	case "clickhouse":
		return "tcp://" + config.Host + ":" + config.Port + "?database=" + config.Dbname + "&username=" + config.User + "&password=" + config.Password + "&read_timeout=10&write_timeout=20"
	default:
		return ""
	}
}
func DsnDecode(dbType, dsn string) (*Config, error) {
	switch dbType {
	case "tidb", "mysql":
		return mysqlDsnDecode(dsn)
	case "pgsql":
		return pgsqlDsnDecode(dsn)
	case "clickhouse":
		return clickhouseDsnDecode(dsn)
	default:
		return nil, errors.New("unsupported database type " + dbType)
	}

}

func mysqlDsnDecode(dsn string) (*Config, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		Type:      "mysql",
		User:      config.User,
		Password:  config.Passwd,
		Dbname:    config.DBName,
		Collation: config.Collation,
		TimeZone:  config.Loc.String(),
	}
	switch config.Net {
	case "tcp":
		conf.Host = config.Addr
	case "unix":
		conf.Host = config.Net + "(" + config.Addr + ")"
	default:
		return nil, errors.New("unsupported network type " + config.Net)
	}
	if config.Net == "unix" {
		conf.Host = config.Addr
	}
	if config.Net == "tcp" {
		addr := strings.Split(config.Addr, ":")
		if len(addr) != 2 {
			return nil, errors.New("invalid address " + config.Addr)
		}
		conf.Host = addr[0]
		conf.Port = addr[1]
	}
	return conf, nil
}

func pgsqlDsnDecode(dsn string) (*Config, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	conf := &Config{
		Type:     "pgsql",
		Host:     config.Host,
		Port:     cast.ToString(config.Port),
		User:     config.User,
		Password: config.Password,
		Dbname:   config.Database,
		DbAlias:  config.Database,
		TimeZone: config.RuntimeParams["timezone"],
	}
	return conf, nil
}

func clickhouseDsnDecode(dsn string) (*Config, error) {
	config, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	if len(config.Addr) == 0 {
		return nil, errors.New("the address cannot be empty")
	}
	addr := strings.Split(config.Addr[0], ":")
	conf := &Config{
		Type:     "clickhouse",
		Host:     addr[0],
		Port:     addr[1],
		User:     config.Auth.Username,
		Password: config.Auth.Password,
		Dbname:   config.Auth.Database,
		DbAlias:  config.Auth.Database,
	}
	return conf, nil
}
