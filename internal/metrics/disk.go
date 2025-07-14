package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MarcusMJV/snapsys/native"
	"golang.org/x/sys/unix"
)

type DiskStats struct {
	Mountpoint string  `json:"mount"`
	UsagePct   float64 `json:"usage_pct"`
	TotalKB    uint64  `json:"total_kb"`
	UsedKB     uint64  `json:"used_kb"`
	FreeKB     uint64  `json:"free_kb"`
}

type DiskMap map[string]native.DiskStats

/*
read_disk_stats
Ignore:
- Virtual filesystems
- Temporary mounts
- Compressed and read-only mounts
- Kernel system paths
*/
var ignoredFSandMounts = map[string]bool{
	"/proc":      true,
	"proc":       true,
	"/sys":       true,
	"sysfs":      true,
	"tmpfs":      true,
	"/run":       true,
	"/snap":      true,
	"devtmpfs":   true,
	"devpts":     true,
	"overlay":    true,
	"squashfs":   true,
	"nsfs":       true,
	"rpc_pipefs": true,
	"cgroup":     true,
	"cgroup2":    true,
	"debugfs":    true,
	"fusectl":    true,
	"tracefs":    true,
}

func GetDiskUsage(mountpoint string) (DiskStats, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(mountpoint, &stat)
	if err != nil {
		return DiskStats{}, fmt.Errorf("statfs failed: %w", err)
	}

	total := stat.Blocks * uint64(stat.Bsize) / 1024
	free := stat.Bfree * uint64(stat.Bsize) / 1024
	used := total - free

	var usage float64
	if total > 0 {
		usage = (float64(used) / float64(total)) * 100
	} else {
		usage = 0.0
	}

	return DiskStats{
		Mountpoint: mountpoint,
		UsagePct:   usage,
		TotalKB:    total,
		UsedKB:     used,
		FreeKB:     free,
	}, nil
}

func GetAllDisks() (DiskMap, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/mounts file: %w", err)
	}
	defer file.Close()

	disks := make(DiskMap)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 3 {
			continue
		}

		mountpoint := fields[1]
		FStype := fields[2]

		if ignoredFSandMounts[mountpoint] || ignoredFSandMounts[FStype] {
			continue
		}

		if _, err := os.Stat(mountpoint); os.IsNotExist(err) {
			continue
		}

		stat, err := native.GetDiskUsage(mountpoint)
		if err == nil {
			if stat.TotalKB != 0 {
				disks[mountpoint] = stat
			}
		}
	}
	return disks, nil
}
