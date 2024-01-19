package manager

import (
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/agclqq/prow-framework/confctr"
)

func TestNew(t *testing.T) {
	type args struct {
		conf Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{conf: Config{Type: CCTypeEtcd, EtcdConf: clientV3.Config{Endpoints: []string{"127.0.0.1:2379"}}}}, wantErr: false},
		{name: "test2", args: args{conf: Config{Type: CCTypeConsul, ConsulConf: &consulapi.Config{Address: "127.0.0.1:8500"}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if _, ok := got.(confctr.CC); !ok {
				t.Errorf("New() got = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
