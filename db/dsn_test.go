package db

import (
	"reflect"
	"testing"
)

func TestDsn(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{config: &Config{
			Type:     "mysql",
			Host:     "mysqlhost.test",
			Port:     "3306",
			User:     "mysqlUser",
			Password: "mysqlPass",
			Dbname:   "mysqldb",
			DbAlias:  "mysqlAlias",
			Charset:  "utf8mb4",
			TimeZone: "utc",
		}}, want: "mysqlUser:mysqlPass@tcp(mysqlhost.test:3306)/mysqldb?charset=utf8mb4&parseTime=True&loc=utc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dsn(tt.args.config); got != tt.want {
				t.Errorf("Dsn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDsnDecode(t *testing.T) {
	type args struct {
		dbType string
		dsn    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{name: "test1", args: args{
			dbType: "mysql",
			dsn:    "mysqlUser:mysqlPass@tcp(mysqlhost.test:3306)/mysqldb?charset=utf8mb4&parseTime=True&loc=utc",
		}, want: &Config{
			Type:      "mysql",
			Host:      "mysqlhost.test",
			Port:      "3306",
			User:      "mysqlUser",
			Password:  "mysqlPass",
			Dbname:    "mysqldb",
			DbAlias:   "mysqldb",
			Charset:   "",
			Collation: "utf8mb4_general_ci",
			TimeZone:  "utc",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DsnDecode(tt.args.dbType, tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("DsnDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DsnDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
