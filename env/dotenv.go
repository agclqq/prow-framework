package env

import (
	"os"
	"strings"

	"github.com/agclqq/godotenv"
)

type DotEnv struct {
	envs map[string]string
}

func NewDotEnv(option ...Option) (Manager, error) {
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
		for k, _ := range m {
			if osV, ok := os.LookupEnv(k); ok {
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
	if e.envs == nil {
		return ""
	}
	return e.envs[key]
}

func (e *DotEnv) SetAll(m map[string]string) bool {
	e.envs = m
	return true
}
func (e *DotEnv) Set(key, val string) bool {
	if e.envs == nil {
		return false
	}
	e.envs[key] = val
	return true
}
