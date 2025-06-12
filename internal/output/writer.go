package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MarcusMJV/snapsys.git/internal/metrics"
)

type Snapshot struct {
	Timestamp time.Time           `json:"timestamp"`
	CPU       metrics.CPUStats    `json:"cpu"`
	Memory    metrics.MemoryStats `json:"memory"`
	Disks     metrics.DiskMap     `json:"disks"`
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
