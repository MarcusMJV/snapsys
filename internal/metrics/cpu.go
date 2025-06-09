package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CPUStats struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
}

func (cpu *CPUStats) Total() uint64 {
	return cpu.User + cpu.Nice + cpu.System + cpu.Idle + cpu.IOWait + cpu.IRQ + cpu.SoftIRQ
}

func (cpu *CPUStats) IdleTime() uint64 {
	return cpu.Idle + cpu.IOWait
}

func CalculateCpuUsage(cpuSnap1, cpuSnap2 CPUStats) float64 {
	deltaTotal := cpuSnap2.Total() - cpuSnap1.Total()
	deltaIdle := cpuSnap2.IdleTime() - cpuSnap1.IdleTime()

	if deltaTotal == 0 {
		return 0.0
	}

	return float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100
}

func ReadCPUStats() (CPUStats, error) {
	var cpuStats CPUStats

	file, err := os.Open("/proc/stat")
	if err != nil {
		return CPUStats{}, err
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
				return CPUStats{}, fmt.Errorf("unecpected format in /proc/stat")
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
