package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/MarcusMJV/snapsys.git/internal/metrics"
)

type Snapshot struct {
	Timestamp time.Time           `json:"timestamp"`
	CPU       metrics.CPUStats    `json:"cpu"`
	Memory    metrics.MemoryStats `json:"memory"`
	Disks     metrics.DiskMap     `json:"disks"`
}

func TakeSnapshot(prevCpuSnap metrics.CPUStatsRaw, outputFile string, now time.Time) bool {
	var wg sync.WaitGroup
	wg.Add(3)
	errChan := make(chan error, 3)
	cpuChan := make(chan metrics.CPUStats, 1)
	memChan := make(chan metrics.MemoryStats, 1)
	diskChan := make(chan metrics.DiskMap, 1)

	go func() {
		defer wg.Done()
		cpu, err := metrics.ReadCPU(&prevCpuSnap)
		if err != nil {
			errChan <- fmt.Errorf("CPU: %w", err)
		}
		cpuChan <- cpu
	}()
	go func() {
		defer wg.Done()
		memory, err := metrics.ReadMemStats()
		if err != nil {
			errChan <- fmt.Errorf("MEMORY: %w", err)
		}
		memChan <- memory
	}()
	go func() {
		defer wg.Done()
		disks, err := metrics.GetAllDisks()
		if err != nil {
			errChan <- fmt.Errorf("DISK: %w", err)
		}
		diskChan <- disks
	}()
	wg.Wait()

	close(errChan)

	var hasError bool
	for err := range errChan {
		fmt.Println(err)
		hasError = true
	}

	if hasError {
		return hasError
	}

	snapshot := Snapshot{
		Timestamp: now,
		CPU:       <-cpuChan,
		Memory:    <-memChan,
		Disks:     <-diskChan,
	}

	err := snapshot.AppendSnapshotJSONL(outputFile)
	if err != nil {
		fmt.Println(err)
		return true
	}

	return false
}

func (s *Snapshot) AppendSnapshotJSONL(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(s); err != nil {
		return fmt.Errorf("flaied to write snaposhot: %w", err)
	}

	return nil
}
