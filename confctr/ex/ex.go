package ex

import (
	consulapi "github.com/hashicorp/consul/api"
	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/agclqq/prow-framework/confctr/manager"
)

func ex() {
	conf := manager.Config{
		Type:       manager.CCTypeConsul,
		EtcdConf:   clientV3.Config{Endpoints: []string{"127.0.0.1:2379"}},
		ConsulConf: &consulapi.Config{Address: "127.0.0.1:8500"},
	}
	cc, err := manager.New(conf)
	if err != nil {
		return
	}
	cc.Get("aa")

}
