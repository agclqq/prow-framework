package args

import (
	"reflect"
	"testing"
)

func Test_tidyParmaWithPrefix(t *testing.T) {
	type args struct {
		param []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{name: "t1", args: args{param: []string{"make:controller", "controllerName", "path", "-r", "y"}}, want: map[string]string{"r": "y"}},
		{name: "t2", args: args{param: []string{"make:controller", "controllerName", "-path", "-r", "y"}}, want: map[string]string{"path": "", "r": "y"}},
		{name: "t3", args: args{param: []string{"make:controller", "-controllerName", "path", "-r", "y"}}, want: map[string]string{"controllerName": "path", "r": "y"}},
		{name: "t4", args: args{param: []string{"make:controller", "-controllerName", "path", "r", "y"}}, want: map[string]string{"controllerName": "path"}},
		{name: "t5", args: args{param: []string{"make:controller", "controllerName", "path", "-r"}}, want: map[string]string{"r": ""}},
		{name: "t6", args: args{param: []string{"make:controller", "controllerName", "-path", "-r"}}, want: map[string]string{"path": "", "r": ""}},
		{name: "t7", args: args{param: []string{}}, want: map[string]string{}},
		{name: "t8", args: args{param: []string{"make:controller", "controllerName", "---", "a"}}, want: map[string]string{}},
		{name: "t9", args: args{param: []string{"make:controller", "controllerName", "--c", "a"}}, want: map[string]string{"c": "a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TidyParmaWithPrefix(tt.args.param); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TidyParmaWithPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
