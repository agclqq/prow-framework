package ca

import (
	"os"
	"testing"
	"time"
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
			ticker := time.NewTicker(3 * time.Second)
			go func() {
				<-ticker.C
				os.Remove("ca.key")
				os.Remove("ca.crt")
				process, err := os.FindProcess(os.Getppid() + 1) //使用debug时，由于debug会创建新的进程，所以要加1
				if err != nil {
					t.Errorf("Svr() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				time.Sleep(2 * time.Second)
				err = process.Signal(os.Interrupt)
				if err != nil {
					t.Errorf("Svr() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}()
			if err := Svr(tt.args.certPath, tt.args.keyPath); (err != nil) != tt.wantErr {
				t.Errorf("Svr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
