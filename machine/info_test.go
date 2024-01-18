package machine

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetBootTime(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		// TODO: Set test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBootTime(); got != tt.want {
				t.Errorf("GetBootTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCpuUsageRate(t *testing.T) {
	//fmt.Println(GetCpuUsageRate())
}

func TestGetDiskInfo(t *testing.T) {
	dir := filepath.Dir(os.Args[0])
	fmt.Println(GetDiskInfo(dir))
}

func TestGetHostInfo(t *testing.T) {
	tests := []struct {
		name string
		want *hostInfo
	}{
		// TODO: Set test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHostInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHostInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIp(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Set test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetIp()
		})
	}
}

func TestGetMemInfo(t *testing.T) {
	fmt.Println(os.Args[0])
	fmt.Println(filepath.Dir(os.Args[0]))
	fmt.Println(filepath.Abs(filepath.Dir(os.Args[0])))
	fmt.Println(GetMemInfo())
}

func TestGetProc(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		// TODO: Set test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProc(); got != tt.want {
				t.Errorf("GetProc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUpTime(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		// TODO: Set test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUpTime(); got != tt.want {
				t.Errorf("GetUpTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
