package env

var DefaultEnvName = "GO_ENV"

type conf struct {
	File    string
	MergeOs bool
	EnvName string
	EnvMaps map[string]string
}
type Option func(em *conf)
type Manager interface {
	GetAll() map[string]string
	SetAll(map[string]string) bool
	Get(key string) string
	Set(key, val string) bool
}
