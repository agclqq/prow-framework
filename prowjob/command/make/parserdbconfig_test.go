package make

import "testing"

func Test_parserDbConfig(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "../../config/db.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parserDbConfig(tt.args.filename)
		})
	}
}
