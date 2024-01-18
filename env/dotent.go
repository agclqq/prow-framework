package env

import (
	"os"
	"strings"

	"github.com/agclqq/godotenv"
)

type DotEnv struct {
	envs map[string]string
}

func NewDotEnv(option ...Option) (*DotEnv, error) {
	cf := &conf{}
	for _, opt := range option {
		opt(cf)
	}
	if cf.File == "" {
		if cf.EnvName == "" {
			cf.EnvName = DefaultEnvName
		}
		cf.File = getEnvFileByOsEnv(cf.EnvName)
	}

	m, err := godotenv.Read(cf.File)
	if err != nil {
		return nil, err
	}
	de := &DotEnv{envs: m}
	if cf.MergeOs {
		for k, v := range m {
			if osV, ok := os.LookupEnv(v); ok {
				de.envs[k] = osV
			}
		}
	}
	return de, nil
}

func getEnvFileByOsEnv(envName string) string {
	var suffix string
	suffix = strings.TrimSpace(os.Getenv(envName))
	file := ".env"
	if "" != suffix {
		file += "." + suffix
	}
	return file
}

func (e *DotEnv) GetAll() map[string]string {
	return e.envs
}
func (e *DotEnv) Get(key string) string {
	return e.envs[key]
}

func (e *DotEnv) SetAll(m map[string]string) bool {
	e.envs = m
	return true
}
func (e *DotEnv) Set(key, val string) bool {
	return e.envs[key] == val
}
