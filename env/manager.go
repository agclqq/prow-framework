package env

import (
	"errors"
)

type Type int

const (
	Dot Type = iota
)

var ErrUnsupportedType = errors.New("unsupported type")
var std, err = New(Dot, WithOsEnv())

// WithOsEnv is a variadic function that accepts a list of environment variables to load.
// If a variable in the system environment is the same as a variable in the configuration file, it is overwritten
func WithOsEnv() Option {
	return func(m *conf) {
		m.MergeOs = true
	}
}

// WithEnvName specifies which environment variable is used to take the value of the current environment
// default is GO_ENV
func WithEnvName(name string) Option {
	return func(m *conf) {
		m.EnvName = name
	}
}

// WithFile is a variadic function that accepts a list of File to load.
// If a variable in the system environment is the same as a variable in the configuration file, it is overwritten
func WithFile(file string) Option {
	return func(m *conf) {
		m.File = file
	}
}

func New(envType Type, opts ...Option) (Manager, error) {
	switch envType {
	case Dot:
		return NewDotEnv(opts...)
	}
	return nil, ErrUnsupportedType
}

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
