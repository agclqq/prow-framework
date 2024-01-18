package env

import (
	"os"
	"strings"
	"testing"
)

var envs = map[string]string{
	"TEST_KEY1": "val",
	"TEST_KEY2": "'val#$%!()2'",
	"TEST_KEY3": "\"val#$%!()3\"",
	"TEST_KEY4": "\"val\\'#\\\"$%!()3\"",
}

func setUp(t *testing.T) {
	f, err := os.Create(".env")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	for k, v := range envs {
		_, err = f.WriteString(k + "=" + v + "\n")
		if err != nil {
			t.Error(err)
			return
		}
	}
}
func tearDown(t *testing.T) {
	err := os.Remove(".env")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNew(t *testing.T) {
	setUp(t)
	defer tearDown(t)
	m, err := New(Dot, WithOsEnv())
	if err != nil {
		t.Error(err)
		return
	}
	for k, v := range envs {
		v = strings.Trim(v, "\"")
		v = strings.Trim(v, "'")
		v = literalConversion(v)

		got := m.Get(k)
		if got != v {
			t.Errorf("want %s, got %s", v, got)
		}
	}
}

func literalConversion(input string) string {
	replacements := map[string]string{
		"\\'":  "'",
		"\\\"": "\"",
		"\\\\": "\\",
	}

	for key, value := range replacements {
		input = strings.Replace(input, key, value, -1)
	}

	return input
}
