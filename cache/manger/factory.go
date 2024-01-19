package manger

import (
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/agclqq/prow-framework/cache"
	predis "github.com/agclqq/prow-framework/cache/redis"
)

type CacheType int

const (
	CacheTypeRedis CacheType = iota
	CacheTypeFile
)

type Config struct {
	Type      CacheType
	RdsConfig *redis.Options
}

func New(conf Config) (cache.Cacher, error) {
	switch conf.Type {
	case CacheTypeRedis:
		return predis.New(conf.RdsConfig), nil
	}
	return nil, errors.New("not support cache type")
}
