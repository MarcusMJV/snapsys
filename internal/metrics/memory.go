package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MemoryStats struct {
	UsagePct float64        `json:"usage_pct"`
	Raw      MemoryStatsRaw `json:"raw"`
}

type MemoryStatsRaw struct {
	MemTotal     uint64 `json:"mem_total_kb"`
	MemUsed      uint64 `json:"mem_used_kb"`
	MemAvailable uint64 `json:"mem_available_kb"`
}

func ReadMemStats() (MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryStats{}, fmt.Errorf("failed to open /proc/meminfo: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mem := make(map[string]uint64)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, _ := strconv.ParseUint(fields[1], 10, 64)

		mem[key] = value
	}

	rawStats := MemoryStatsRaw{
		MemAvailable: mem["MemAvailable"],
		MemTotal:     mem["MemTotal"],
		MemUsed:      mem["MemTotal"] - mem["MemAvailable"],
	}

	usage := (float64(rawStats.MemUsed) / float64(rawStats.MemTotal)) * 100
	return MemoryStats{UsagePct: usage, Raw: rawStats}, nil
}
