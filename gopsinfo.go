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

type psInfo struct {
	dateTime        string
	logicalCores    int
	physicalCores   int
	percentPerCpu   []float64
	cpuPercent      float64
	cpuModel        []string
	memTotal        uint64
	memUsed         uint64
	memUsedPercent  float64
	bytesRecv       float64
	bytesSent       float64
	diskTotal       uint64
	diskUsed        uint64
	diskUsedPercent float64
	load            []string
	os              string
	platform        string
	platformFamily  string
	platformVersion string
}

var (
	timeFormat = "2006-01-02T15:04:05"
	recv       float64
	sent       float64
	pi         psInfo
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

	pi.memTotal = v.Total
	pi.logicalCores = logicalCount
	pi.physicalCores = physicalCount
	pi.cpuModel = cpuModel
	pi.os = hostInfoStat.OS
	pi.platform = hostInfoStat.Platform
	pi.platformFamily = hostInfoStat.PlatformFamily
	pi.platformVersion = hostInfoStat.PlatformVersion
}

func GetPsInfo(interval int) psInfo {
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

	pi.diskTotal = diskTotal
	pi.diskUsedPercent = diskUsedPercent * 100
	pi.diskUsed = diskUsed
	pi.memUsedPercent = v.UsedPercent
	pi.memUsed = v.Used
	pi.percentPerCpu = percentPerCpu
	pi.cpuPercent = cpuPercent[0]
	pi.bytesRecv = recvRate
	pi.bytesSent = sentRate
	pi.dateTime = time.Now().Format(timeFormat)
	pi.load = []string{
		parseFloatNum(loadAvg.Load1),
		parseFloatNum(loadAvg.Load5),
		parseFloatNum(loadAvg.Load15),
	}

	return pi
}

func parseFloatNum(n float64) string {
	return fmt.Sprintf("%.2f", n)
}

func errHandler(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
