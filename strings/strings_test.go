package strings

import "testing"

func TestToLowFirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "l1", args: args{str: "ABC0"}, want: "aBC0"},
		{name: "l2", args: args{str: "aBC0"}, want: "aBC0"},
		{name: "l3", args: args{str: "abC0"}, want: "abC0"},
		{name: "l4", args: args{str: "0BC0"}, want: "0BC0"},
		{name: "l5", args: args{str: "工BC0"}, want: "工BC0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLowFirst(tt.args.str); got != tt.want {
				t.Errorf("ToLowFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUpFirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "u1", args: args{str: "ABC0"}, want: "ABC0"},
		{name: "u2", args: args{str: "aBC0"}, want: "ABC0"},
		{name: "u3", args: args{str: "abC0"}, want: "AbC0"},
		{name: "u4", args: args{str: "0BC0"}, want: "0BC0"},
		{name: "u5", args: args{str: "工BC0"}, want: "工BC0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToUpFirst(tt.args.str); got != tt.want {
				t.Errorf("ToUpFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
