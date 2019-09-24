# gopsinfo

## 简介
获取机器基础数据，包括CPU，磁盘，网络，内存相关信息，按配置写入指定日志文件，供下游分析展现等。

## 使用
```
import "github.com/schoeu/gopsinfo"
```

### 方法

```
// GetPsInfo(during time.Duration) during为检测时间区间
gopsinfo.GetPsInfo(during)
```

## 数据
获取到的字段

|占位符|含义|示例|备注|
|--|--|--|--|
|dateTime|日期时间戳|2019-06-28T17:37:11|当前时间戳|
|logicalCores|逻辑核数|8||
|physicalCores|物理核数|4||
|percentPerCpu|单cpu使用率|[33.66 3.00 30.00 3.96 30.00 3.96 27.72 3.96]|展现每一个逻辑核的使用率|
|cpuPercent|cpu综合使用率|6.64|使用率为6.64%|
|cpuModel|cpu型号|"Intel(R) Core(TM) i7-4750HQ CPU @ 2.00GHz"|多类核会以`,`隔开|
|memTotal|总内存|8192MB|8GB，此处以MB来展现|
|memUsed|已使用内存|5516.53MB|已使用了5516.53MB|
|memUsedPercent|内存使用率|67.34|已使用占比67.34%|
|bytesRecv|网卡下行速率|4.00KB/s|下行速率|
|bytesSent|网卡上行速率|1.50KB/s|上行速率|
|diskTotal|磁盘总空间|467GB|磁盘总计467G，不包括隐藏分区|
|diskUsed|磁盘已使用空间|159GB|已使用159GB|
|diskUsedPercent|磁盘使用占比|34.20|磁盘使用了34.20%|
|os|系统类型|darwin||
|platform|系统所属平台|darwin|
|platformFamily|系统平台分类| Standalone Workstation|
|platformVersion|系统版本|10.14.5|

## MIT License

Copyright (c) 2019 Schoeu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
