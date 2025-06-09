package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	CPUUsage  float64   `json:"cpu_usage"`
}

func WriteSnapshotsToFile(snapshots []Snapshot, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(snapshots); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
