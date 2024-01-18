package module

import "testing"

func TestGetModuleName(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"test", "github.com/agclqq/prow-framework", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetModuleName()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetModuleName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetModuleName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
