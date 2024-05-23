package env

import "testing"

func TestDotEnv_GetNoneEnvs(t *testing.T) {
	dot := &DotEnv{}
	if dot.Get("key") != "" {
		t.Errorf("want %s, got %s", "", dot.Get("key"))
	}
	if dot.Set("key", "val") {
		t.Error("want false, got true")
	}
}
