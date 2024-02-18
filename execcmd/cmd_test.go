package execcmd

import (
	"reflect"
	"testing"
)

func TestCommand(t *testing.T) {
	type args struct {
		name string
		arg  []string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "t1", args: args{name: "notexist", arg: []string{"-h"}}, wantErr: true, want: nil},
		{name: "t2", args: args{name: "echo", arg: []string{"hello"}}, wantErr: false, want: []byte("hello\n")},
		{name: "t3", args: args{name: "ls", arg: []string{"-"}}, wantErr: true, want: []byte("ls: -: No such file or directory\n")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Command(tt.args.name, tt.args.arg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Command() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command() got = %v, want %v", got, tt.want)
			}
		})
	}
}
