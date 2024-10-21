package autoenv

import "github.com/agclqq/prow-framework/env"

var std, err = env.New(env.Dot, env.WithOsEnv())

func Get(key string, defaultVal ...string) string {
	if std == nil {
		return ""
	}
	if v := std.Get(key); v != "" {
		return v
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""
}
func GetAll() map[string]string {
	if std == nil {
		return nil
	}
	return std.GetAll()
}
