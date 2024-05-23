package disk

import "testing"

func TestBlock(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "t1", want: 4096},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Block(); got != tt.want {
				t.Errorf("Block() = %v, want %v", got, tt.want)
			}
		})
	}
}
