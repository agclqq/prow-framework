package ca

import (
	"os"
	"testing"
)

func TestSvr(t *testing.T) {
	defer func() {
		os.Remove("ca.key")
		os.Remove("ca.crt")
	}()
	type args struct {
		keyPath  string
		certPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "t1", args: args{keyPath: "ca.key", certPath: "ca.crt"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Svr(tt.args.certPath, tt.args.keyPath); (err != nil) != tt.wantErr {
				t.Errorf("Svr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
