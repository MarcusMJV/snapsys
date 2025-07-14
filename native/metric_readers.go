package native

/*
#include <stdlib.h>
#include "metric_readers.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type CPUStatsRaw struct {
	User    uint64 `json:"user"`
	Nice    uint64 `json:"nice"`
	System  uint64 `json:"system"`
	Idle    uint64 `json:"idle"`
	IOWait  uint64 `json:"iowait"`
	IRQ     uint64 `json:"irq"`
	SoftIRQ uint64 `json:"softirq"`
}

type MemoryStatsRaw struct {
	MemTotal     uint64 `json:"mem_total"`
	MemFree      uint64 `json:"mem_free"`
	MemAvailable uint64 `json:"mem_available"`
	Buffers      uint64 `json:"buffers"`
	Cached       uint64 `json:"cached"`
}

type DiskStats struct {
	UsagePct float32 `json:"usage_pct"`
	TotalKB  uint64  `json:"total_kb"`
	UsedKB   uint64  `json:"used_kb"`
	FreeKB   uint64  `json:"free_kb"`
}

func (cpu *CPUStatsRaw) Total() uint64 {
	return cpu.User + cpu.Nice + cpu.System + cpu.Idle + cpu.IOWait + cpu.IRQ + cpu.SoftIRQ
}

func (cpu *CPUStatsRaw) IdleTime() uint64 {
	return cpu.Idle + cpu.IOWait
}

func ReadCPUStatsRawC() CPUStatsRaw {
	raw := C.read_proc_stat()

	return CPUStatsRaw{
		User:    uint64(raw.user),
		Nice:    uint64(raw.nice),
		System:  uint64(raw.system),
		Idle:    uint64(raw.idle),
		IOWait:  uint64(raw.iowait),
		IRQ:     uint64(raw.irq),
		SoftIRQ: uint64(raw.softirq),
	}
}

func ReadMemStatsRawC() MemoryStatsRaw {
	raw := C.read_proc_meminfo()

	return MemoryStatsRaw{
		MemTotal:     uint64(raw.mem_total),
		MemFree:      uint64(raw.mem_free),
		MemAvailable: uint64(raw.mem_available),
		Buffers:      uint64(raw.buffers),
		Cached:       uint64(raw.cached),
	}

}

func GetDiskUsage(mountpoint string) (DiskStats, error) {
	cstr := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(cstr))

	var result C.DiskStats
	status := C.read_disk_stats(cstr, &result)
	if status != 0 || result.success == 0 {
		return DiskStats{}, fmt.Errorf("statfs failed fot %s", mountpoint)
	}

	return DiskStats{
		TotalKB:  uint64(result.total_kb),
		UsedKB:   uint64(result.used_kb),
		FreeKB:   uint64(result.free_kb),
		UsagePct: float32(result.usage_pct),
	}, nil
}
