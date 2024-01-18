package prettytable

import (
	"testing"
)

func Test_printTextTable(t *testing.T) {
	type args struct {
		data [][]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{data: [][]string{{"Header 1", "Header 2", "Header 3"}, {"Data 1", "Data 2\nMultiline", "Data 3"}}}, want: ` ________ ________________ ________
|Header 1|Header 2        |Header 3|
|Data 1  |Data 2 Multiline|Data 3  |
 ‾‾‾‾‾‾‾‾ ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾ ‾‾‾‾‾‾‾‾
`},
		{name: "test2", args: args{data: [][]string{{}}}, want: `
|

`},
		{name: "test3", args: args{data: [][]string{}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PlainText(tt.args.data)
			if got != tt.want {
				t.Errorf("printTextTable() =\n%v, want\n%v", got, tt.want)
			}
		})
	}
}

func Test_markdown(t *testing.T) {
	type args struct {
		data [][]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{data: [][]string{{"Header 1", "Header 2", "Header 3"}, {"Data 1", "Data 2\nMultiline", "Data 3"}}}, want: `| Header 1 | Header 2 | Header 3 |
| Data 1 | Data 2 Multiline | Data 3 |
`},
		{name: "test2", args: args{data: [][]string{{}}}, want: `|
`},
		{name: "test3", args: args{data: [][]string{}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Markdown(tt.args.data); got != tt.want {
				t.Errorf("Markdown() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
