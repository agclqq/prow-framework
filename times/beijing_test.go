package times

import (
	"fmt"
	"testing"
)

func TestStringToTime(t *testing.T) {
	layout := FormatDate + " " + FormatShortTime
	s1 := "2022-06-22 20:20"
	t1 := StringToTime(layout, s1)
	rs1 := TimeToString(layout, t1)
	if s1 != rs1 {
		t.Errorf("want %s ,got %s", s1, rs1)
	}
	fmt.Println(s1, rs1)
}
