package gopsinfo

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type PsInfo struct {
	DateTime        string
	LogicalCores    int
	PhysicalCores   int
	PercentPerCpu   []float64
	CpuPercent      float64
	CpuModel        []string
	MemTotal        uint64
	MemUsed         uint64
	MemUsedPercent  float64
	RecvRate        float64
	SentRate        float64
	DiskTotal       uint64
	DiskUsed        uint64
	DiskUsedPercent float64
	Load            []string
	Os              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
}

var (
	timeFormat = "2006-01-02T15:04:05"
	recv       float64
	sent       float64
	pi         PsInfo
)

func init() {
	getSysInfo()
}

func getSysInfo() {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	logicalCount, _ := cpu.Counts(true)
	physicalCount, _ := cpu.Counts(false)

	var cpuModel []string
	for _, subCpu := range c {
		cpuModel = append(cpuModel, fmt.Sprintf(`"%s"`, subCpu.ModelName))
	}

	hostInfoStat, _ := host.Info()

	pi.MemTotal = v.Total
	pi.LogicalCores = logicalCount
	pi.PhysicalCores = physicalCount
	pi.CpuModel = cpuModel
	pi.Os = hostInfoStat.OS
	pi.Platform = hostInfoStat.Platform
	pi.PlatformFamily = hostInfoStat.PlatformFamily
	pi.PlatformVersion = hostInfoStat.PlatformVersion
}

func GetPsInfo(interval int) PsInfo {
	v, _ := mem.VirtualMemory()
	percentPerCpu, _ := cpu.Percent(time.Second, true)
	cpuPercent, _ := cpu.Percent(time.Second, false)
	diskInfo, _ := disk.Partitions(true)
	loadAvg, _ := load.Avg()
	var diskTotal, diskUsed uint64
	for _, v := range diskInfo {
		device := v.Mountpoint
		distDetial, _ := disk.Usage(device)
		if distDetial != nil {
			diskTotal += distDetial.Total
			diskUsed += distDetial.Used
		}
	}

	nw, _ := net.IOCounters(false)
	parseNum := float64(uint64(interval) / 1000)
	var recvRate, sentRate float64
	if len(nw) > 0 && nw[0].Name == "all" {
		br := float64(nw[0].BytesRecv)
		bs := float64(nw[0].BytesSent)
		recvRate = (br - recv) / parseNum
		sentRate = (bs - sent) / parseNum

		// 初次获取上下行信息，矫正数据
		if recv == 0 || sent == 0 {
			recvRate = 0
			sentRate = 0
		}

		recv = br
		sent = bs
	}

	diskUsedPercent := float64(diskUsed) / float64(diskTotal)

	pi.DiskTotal = diskTotal
	pi.DiskUsedPercent = diskUsedPercent * 100
	pi.DiskUsed = diskUsed
	pi.MemUsedPercent = v.UsedPercent
	pi.MemUsed = v.Used
	pi.PercentPerCpu = percentPerCpu
	pi.CpuPercent = cpuPercent[0]
	pi.RecvRate = recvRate
	pi.SentRate = sentRate
	pi.DateTime = time.Now().Format(timeFormat)
	pi.Load = []string{
		parseFloatNum(loadAvg.Load1),
		parseFloatNum(loadAvg.Load5),
		parseFloatNum(loadAvg.Load15),
	}

	return pi
}

func parseFloatNum(n float64) string {
	return fmt.Sprintf("%.2f", n)
}
