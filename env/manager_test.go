package env

import (
	"os"
	"os/exec"
	"testing"
)

var envs = map[string]string{
	"TEST_KEY1": "=val",
	"TEST_KEY2": "'val=#$%!()2'",
	"TEST_KEY3": "\"val#=$%!()3\"",
	"TEST_KEY4": "\"val='#\\\"$%!()3\"",
	"TEST_KEY5": "test os env",
	"TEST_KEY6": "",
}
var envRs = map[string]string{
	"TEST_KEY1": "=val",
	"TEST_KEY2": "val=#$%!()2",
	"TEST_KEY3": "val#=$%!()3",
	"TEST_KEY4": "val='#\"$%!()3",
	"TEST_KEY5": "new os env", //在程序中已通过OS环境变量修改
	"TEST_KEY6": "default",
}

func setUp(t *testing.T, envName string) {
	f, err := os.Create(envName)
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
	os.Setenv("TEST_KEY5", "new os env")
	// 读取文件内容
	cmd := exec.Command("cat", envName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
func tearDown(t *testing.T, envName string) {
	err := os.Remove(envName)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUnsupported(t *testing.T) {
	_, err := New(999)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestNoFile(t *testing.T) {
	_, err := New(Dot, WithOsEnv())
	if err == nil {
		t.Error("want error, got nil")
	}
	if Get("not_exist") != "" {
		t.Error("want empty, got not empty")
	}
}
func TestNew(t *testing.T) {
	envName := ".env"
	setUp(t, envName)
	defer tearDown(t, envName)

	m, err := New(Dot, WithOsEnv())
	if err != nil {
		t.Error(err)
		return
	}
	baseVerify(t, m)
}

func TestSet(t *testing.T) {
	envName := ".env"
	setUp(t, envName)
	defer tearDown(t, envName)

	m, err := New(Dot, WithOsEnv())
	if err != nil {
		t.Error(err)
		return
	}
	newEnv := map[string]string{
		"TEST_KEY1": "new val 1",
		"TEST_KEY2": "new val 2",
		"TEST_KEY3": "new val 3",
		"TEST_KEY4": "new val 4",
		"TEST_KEY5": "new val 5",
	}
	m.SetAll(newEnv)
	m.Set("TEST_KEY1", "set1")
	if m.Get("TEST_KEY1") != "set1" {
		t.Errorf("want %s, got %s", "set1", m.Get("TEST_KEY1"))
	}
	if m.Get("TEST_KEY2") != newEnv["TEST_KEY2"] {
		t.Errorf("want %s, got %s", newEnv["TEST_KEY2"], m.Get("TEST_KEY2"))
	}
	if m.Get("TEST_KEY3") != newEnv["TEST_KEY3"] {
		t.Errorf("want %s, got %s", newEnv["TEST_KEY3"], m.Get("TEST_KEY3"))
	}
	if m.Get("TEST_KEY4") != newEnv["TEST_KEY4"] {
		t.Errorf("want %s, got %s", newEnv["TEST_KEY4"], m.Get("TEST_KEY4"))
	}
	if m.Get("TEST_KEY5") != newEnv["TEST_KEY5"] {
		t.Errorf("want %s, got %s", newEnv["TEST_KEY5"], m.Get("TEST_KEY5"))
	}
}

func TestProd(t *testing.T) {
	envName := ".env.prod"
	setUp(t, envName)
	defer tearDown(t, envName)

	os.Setenv("GOGO_ENV", "prod")
	m, err := New(Dot, WithOsEnv(), WithEnvName("GOGO_ENV"))
	if err != nil {
		t.Error(err)
		return
	}
	baseVerify(t, m)
}

func TestDiyFile(t *testing.T) {
	envName := ".diy"
	setUp(t, envName)
	defer tearDown(t, envName)

	m, err := New(Dot, WithOsEnv(), WithFile(envName))
	if err != nil {
		t.Error(err)
		return
	}
	baseVerify(t, m)
}

func TestStd(t *testing.T) {
	envName := ".env"
	setUp(t, envName)
	defer tearDown(t, envName)

	std, err = New(Dot, WithOsEnv())

	baseVerify(t, nil)

	rs := GetAll()
	if rs["TEST_KEY1"] != envRs["TEST_KEY1"] {
		t.Errorf("want %s, got %s", envRs["TEST_KEY1"], rs["TEST_KEY1"])
	}
	if rs["TEST_KEY2"] != envRs["TEST_KEY2"] {
		t.Errorf("want %s, got %s", envRs["TEST_KEY2"], rs["TEST_KEY2"])
	}
	if rs["TEST_KEY3"] != envRs["TEST_KEY3"] {
		t.Errorf("want %s, got %s", envRs["TEST_KEY3"], rs["TEST_KEY3"])
	}
	if rs["TEST_KEY4"] != envRs["TEST_KEY4"] {
		t.Errorf("want %s, got %s", envRs["TEST_KEY4"], rs["TEST_KEY4"])
	}
	if rs["TEST_KEY5"] != envRs["TEST_KEY5"] {
		t.Errorf("want %s, got %s", envRs["TEST_KEY5"], rs["TEST_KEY5"])
	}
	if Get("TEST_KEY6", "default") != "default" {
		t.Errorf("want %s, got %s", "default", Get("TEST_KEY	6", "default"))
	}
	if Get("NOT_EXIST_KEY") != "" {
		t.Errorf("want %s, got %s", "", Get("NOT_EXIST_KEY"))
	}
}

func baseVerify(t *testing.T, m Manager) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Example 1", args: args{"TEST_KEY1"}, want: envRs["TEST_KEY1"]},
		{name: "Example 2", args: args{"TEST_KEY2"}, want: envRs["TEST_KEY2"]},
		{name: "Example 3", args: args{"TEST_KEY3"}, want: envRs["TEST_KEY3"]},
		{name: "Example 4", args: args{"TEST_KEY4"}, want: envRs["TEST_KEY4"]},
		{name: "Example 5", args: args{"TEST_KEY5"}, want: envRs["TEST_KEY5"]},
	}
	got := ""
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if m == nil {
				got = Get(tt.args.key)
			} else {
				got = m.Get(tt.args.key)
			}
			if got != tt.want {
				t.Errorf("want %s, got %s", tt.want, got)
			}
		})
	}
}
