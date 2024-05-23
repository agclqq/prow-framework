package replaceholder

import (
	"encoding/json"
	"testing"
)

func Test_strReplace(t *testing.T) {
	m := map[string]interface{}{}
	type args struct {
		s   string
		old string
		new string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "t1", args: args{s: `{"name":"aa{{test}}aa"}`, old: "{{test}}", new: "xiaoming"}, want: `{"name":"aaxiaomingaa"}`, wantErr: false},
		{name: "t2", args: args{s: `{"age":{{age}}}`, old: "{{age}}", new: "18"}, want: `{"age":18}`, wantErr: false},
		{name: "t3", args: args{s: `{"other":{{other}}}`, old: "{{other}}", new: `["aa","bb"]`}, want: `{"other":["aa","bb"]}`, wantErr: false},
		//{name: "t4", args: args{s: `{"name":"{{test}}","age:{{age}},"other":{{other}}}`, old: "{{other}}", new: `{"aa":1,"bb":"b","cc":["c",{"a":"a"}]}`}, want: `{"name":"{{test}}","age:{{age}},"other":{"aa":1,"bb":"b","cc":["c",{"a":"a"}]}}`, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strReplace(tt.args.s, tt.args.old, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("strReplace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("strReplace() got = %v, want %v", got, tt.want)
			}
			err = json.Unmarshal([]byte(got), &m)
			if err != nil {
				t.Error(err)
			}

		})
	}
}
