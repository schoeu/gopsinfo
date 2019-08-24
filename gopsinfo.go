package gopsinfo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type PsInfo struct {
	DateTime        string  `json:"dataTime"`
	LogicalCores    int     `json:"logicalCores"`
	PhysicalCores   int     `json:"physicalCores"`
	PercentPerCpu   string  `json:"percentPerCpu"`
	CpuPercent      float64 `json:"cpuPercent"`
	CpuModel        string  `json:"cpuModel"`
	MemTotal        uint64  `json:"memTotal"`
	MemUsed         uint64  `json:"memUsed"`
	MemUsedPercent  float64 `json:"memUsedPercent"`
	RecvSpeed       float64 `json:"recvSpeed"`
	SentSpeed       float64 `json:"sentSpeed"`
	DiskTotal       uint64  `json:"diskTotal"`
	DiskUsed        uint64  `json:"diskUsed"`
	DiskUsedPercent float64 `json:"diskUsedPercent"`
	Load            string  `json:"load"`
	Os              string  `json:"os"`
	Platform        string  `json:"platform"`
	PlatformFamily  string  `json:"platformFamily"`
	PlatformVersion string  `json:"platformVersion"`
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
		modeName := strings.Replace(subCpu.ModelName, "\"", "", -1)
		cpuModel = append(cpuModel, fmt.Sprintf(`"%s"`, modeName))
	}

	hostInfoStat, _ := host.Info()

	pi.MemTotal = v.Total
	pi.LogicalCores = logicalCount
	pi.PhysicalCores = physicalCount
	pi.CpuModel = strings.Join(cpuModel, ",")
	pi.Os = hostInfoStat.OS
	pi.Platform = hostInfoStat.Platform
	pi.PlatformFamily = hostInfoStat.PlatformFamily
	pi.PlatformVersion = hostInfoStat.PlatformVersion
}

func GetPsInfo(interval int) PsInfo {
	v, _ := mem.VirtualMemory()
	percentPerCpu, _ := cpu.Percent(time.Microsecond, true)

	var perCpuData []string
	for _, v := range percentPerCpu {
		perCpuData = append(perCpuData, strconv.FormatFloat(v, 'f', 2, 64))
	}

	cpuPercent, _ := cpu.Percent(time.Microsecond, false)
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
	parseNum := float64(interval / 1000)
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
	pi.PercentPerCpu = strings.Join(perCpuData, ",")
	pi.CpuPercent = cpuPercent[0]
	pi.RecvSpeed = recvRate
	pi.SentSpeed = sentRate
	pi.DateTime = time.Now().Format(timeFormat)
	pi.Load = strings.Join([]string{
		parseFloatNum(loadAvg.Load1),
		parseFloatNum(loadAvg.Load5),
		parseFloatNum(loadAvg.Load15),
	}, ",")

	return pi
}

func parseFloatNum(n float64) string {
	return fmt.Sprintf("%.2f", n)
}
