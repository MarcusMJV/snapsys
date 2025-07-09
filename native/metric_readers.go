package native

/*
#include "metric_readers.h"
*/
import "C"

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
