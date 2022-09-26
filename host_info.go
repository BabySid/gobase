package gobase

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"net"
	"time"
)

func GetCPUUsage() float64 {
	percents, _ := cpu.Percent(time.Second*1, false)
	return percents[0]
}

func GetMEMUsage() float64 {
	mem, _ := mem.VirtualMemory()
	return mem.UsedPercent
}

func GetUpTime() uint64 {
	uptime, _ := host.Uptime()
	return uptime
}

func GetDiskPartitions() []string {
	ps, _ := disk.Partitions(false)
	rs := make([]string, len(ps), len(ps))
	for i := range ps {
		rs[i] = ps[i].Device
	}

	return rs
}

func GetLocalIP() (string, error) {
	iFaces, e := net.Interfaces()
	if e != nil {
		return "", e
	}

	for _, face := range iFaces {
		addrs, err := face.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if ip.IsGlobalUnicast() {
					return ip.String(), nil
				}
			}
		}
	}

	return "", nil
}
