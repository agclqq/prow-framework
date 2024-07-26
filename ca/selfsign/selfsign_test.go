package selfsign

import (
	"fmt"
	"testing"
)

func TestCa_Sign(t *testing.T) {
	tests := []struct {
		name    string
		ca      *Ca
		wantErr bool
	}{
		{name: "t1", ca: NewCa([]string{"CN"}, []string{"Beijing"}, []string{"Beijing"}, []string{"my_company"}, []string{"my_department"}, "my_name"), wantErr: false},
		{name: "t2", ca: NewCa([]string{"CN"}, []string{"Beijing"}, []string{"Beijing"}, []string{"my_company"}, []string{"my_department"}, "my_name", WithBit(1024)), wantErr: false},
		{name: "t3", ca: NewCa([]string{"CN"}, []string{"Beijing"}, []string{"Beijing"}, []string{"my_company"}, []string{"my_department"}, "my_name", WithDays(100)), wantErr: false},
		{name: "t4", ca: NewCa([]string{"CN"}, []string{"Beijing"}, []string{"Beijing"}, []string{"my_company"}, []string{"my_department"}, "my_name", WithDays(100), WithBit(2048)), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cert, key, err := tt.ca.Sign()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(string(key), string(cert))
		})
	}
}
