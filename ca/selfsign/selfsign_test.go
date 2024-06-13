package selfsign

import (
	"fmt"
	"os/user"
	"path/filepath"
	"testing"
)

var homeDir = ""

func homePath() (string, error) {
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// 获取家目录
	homeDir = currentUser.HomeDir
	return homeDir, nil
}

func TestNewCa(t *testing.T) {
	_, err := homePath()
	if err != nil {
		t.Errorf("homePath() error = %v", err)
		return
	}
	type args struct {
		c    string
		st   string
		l    string
		o    string
		ou   string
		cn   string
		opts []CaOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "t1", args: args{c: "CN", st: "Beijing", l: "Beijing", o: "my_company", ou: "my_department", cn: "my_name"}, wantErr: false},
		{name: "t2", args: args{c: "CN", st: "Beijing", l: "Beijing", o: "my_company", ou: "my_department", cn: "my_name", opts: []CaOption{WithDir(filepath.Join(homeDir, "/ca/root2")), WithDays(100)}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privetKey, pem, err := NewCa(tt.args.c, tt.args.st, tt.args.l, tt.args.o, tt.args.ou, tt.args.cn, tt.args.opts...).Sign()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCa() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("privetKey: %s ; pem: %s", privetKey, pem)
		})
	}
}
