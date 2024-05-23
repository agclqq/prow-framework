package manager

import (
	"errors"

	consulapi "github.com/hashicorp/consul/api"
	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/agclqq/prow-framework/confctr"
	"github.com/agclqq/prow-framework/confctr/consul"
	"github.com/agclqq/prow-framework/confctr/etcd"
)

type CacheType int

const (
	CCTypeEtcd CacheType = iota
	CCTypeConsul
)

type Config struct {
	Type       CacheType
	EtcdConf   clientV3.Config
	ConsulConf *consulapi.Config
}

func New(conf Config) (confctr.CC, error) {
	switch conf.Type {
	case CCTypeEtcd:
		return etcd.New(conf.EtcdConf)
	case CCTypeConsul:
		return consul.New(conf.ConsulConf)
	}
	return nil, errors.New("not support cache type")
}
