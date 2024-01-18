package machine

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

func GetIp() {

}

type hostInfo struct {
	Hostname             string `json:"hostname"`
	OS                   string `json:"os"`              // ex: freebsd, linux
	Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily       string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion      string `json:"platformVersion"` // version of the complete OS
	KernelVersion        string `json:"kernelVersion"`   // version of the OS kernel (if available)
	KernelArch           string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string `json:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole"` // guest or host
	HostID               string `json:"hostId"`             // ex: uuid
}

var hostV *hostInfo

//---host------------

func GetHostInfo() *hostInfo {
	if hostV == nil {
		v, err := host.Info()
		if err == nil {
			hostV = &hostInfo{
				Hostname:             v.Hostname,
				OS:                   v.OS,
				Platform:             v.Platform,
				PlatformFamily:       v.PlatformFamily,
				PlatformVersion:      v.PlatformVersion,
				KernelVersion:        v.KernelVersion,
				KernelArch:           v.KernelArch,
				VirtualizationSystem: v.VirtualizationSystem,
				VirtualizationRole:   v.VirtualizationRole,
				HostID:               v.HostID,
			}
		}

	}
	return hostV
}

func GetUpTime() uint64 {
	v, _ := host.Uptime()
	return v
}
func GetBootTime() uint64 {
	v, _ := host.BootTime()
	return v
}
func GetProc() uint64 {
	procs, _ := process.Pids()
	return uint64(len(procs))
}

func GetCpuUsageRate(interval time.Duration) float64 {
	v, _ := cpu.Percent(interval, false)
	if len(v) >= 1 {
		return v[0]
	}
	return 0
}

func GetMemInfo() *mem.VirtualMemoryStat {
	v, _ := mem.VirtualMemory()
	return v
}

func GetDiskInfo(path string) *disk.UsageStat {
	v, _ := disk.Usage(path)
	return v
}
