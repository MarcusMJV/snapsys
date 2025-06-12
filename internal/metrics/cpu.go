package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CPUStats struct {
	UsagePct float64     `json:"usage_pct"`
	Raw      CPUStatsRaw `json:"raw"`
}

type CPUStatsRaw struct {
	User    uint64 `json:"user"`
	Nice    uint64 `json:"nice"`
	System  uint64 `json:"system"`
	Idle    uint64 `json:"idle"`
	IOWait  uint64 `json:"iowait"`
	IRQ     uint64 `json:"irq"`
	SoftIRQ uint64 `json:"softirq"`
}

func ReadCPU(prevCpuSnap *CPUStatsRaw) (CPUStats, error) {
	cpuRaw, err := ReadCPUStatsRaw()
	if err != nil {
		return CPUStats{}, err
	}

	cpuUsage := CalculateCpuUsage(prevCpuSnap, &cpuRaw)

	prevCpuSnap = &cpuRaw
	return CPUStats{UsagePct: cpuUsage, Raw: cpuRaw}, nil
}

func (cpu *CPUStatsRaw) Total() uint64 {
	return cpu.User + cpu.Nice + cpu.System + cpu.Idle + cpu.IOWait + cpu.IRQ + cpu.SoftIRQ
}

func (cpu *CPUStatsRaw) IdleTime() uint64 {
	return cpu.Idle + cpu.IOWait
}

func CalculateCpuUsage(cpuSnap1, cpuSnap2 *CPUStatsRaw) float64 {
	deltaTotal := cpuSnap2.Total() - cpuSnap1.Total()
	deltaIdle := cpuSnap2.IdleTime() - cpuSnap1.IdleTime()

	if deltaTotal == 0 {
		return 0.0
	}

	return float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100
}

func ReadCPUStatsRaw() (CPUStatsRaw, error) {
	var cpuStats CPUStatsRaw

	file, err := os.Open("/proc/stat")
	if err != nil {
		return CPUStatsRaw{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		/* Ignoring cpu0, cpu1, cpu2...etc.
		Foucusing on total CPU stats for now*/
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)

			if len(fields) < 8 {
				return CPUStatsRaw{}, fmt.Errorf("unecpected format in /proc/stat")
			}

			cpuStats.User, _ = strconv.ParseUint(fields[1], 10, 64)
			cpuStats.Nice, _ = strconv.ParseUint(fields[2], 10, 64)
			cpuStats.System, _ = strconv.ParseUint(fields[3], 10, 64)
			cpuStats.Idle, _ = strconv.ParseUint(fields[4], 10, 64)
			cpuStats.IOWait, _ = strconv.ParseUint(fields[5], 10, 64)
			cpuStats.IRQ, _ = strconv.ParseUint(fields[6], 10, 64)
			cpuStats.SoftIRQ, _ = strconv.ParseUint(fields[7], 10, 64)

			break
		}
	}

	return cpuStats, nil
}
