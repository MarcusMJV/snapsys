/*
SnapSys - Lightweight System Benchmarking Tool
Copyright © 2025 Marcus Vorster
Released under the MIT License
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MarcusMJV/snapsys/internal/metrics"
	"github.com/MarcusMJV/snapsys/internal/output"
	"github.com/spf13/cobra"
)

type Snapshot struct {
}

var interval time.Duration
var duration time.Duration
var outputFile string

const (
	green = "\033[32m"
	blue  = "\033[34m"
	red   = "\033[31m"
	reset = "\033[0m"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Capture system metrics over a time period",
	Long: `The 'snapshot' command collects CPU, memory, and disk usage statistics 
at fixed intervals over a specified duration.

Examples:
  snapsys snapshot --duration 30s --interval 5s --output stats.jsonl

This will record system metrics every 5 seconds for 30 seconds and save 
them to 'stats.jsonl' in JSONL format.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting snapshot run...")
		runSnapshot()
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().DurationVar(&interval, "interval", 3*time.Second, "Interval between snapshots")
	snapshotCmd.Flags().DurationVar(&duration, "duration", 30*time.Second, "Total duration to run snapshots")
	snapshotCmd.Flags().StringVar(&outputFile, "output", "", "Output file path (defaults to snapsys_runs/snaprun_TIMESTAMP.jsonl)")
}

func runSnapshot() {
	if interval < time.Second {
		fmt.Println(green + "Minimum supported interval is 1s. Using 1s instead." + reset)
		interval = time.Second
	}
	if outputFile == "" {
		defaultPath, err := getDefaultOutputPath()
		if err != nil {
			fmt.Println("Failed to determine output file path:", err)
			os.Exit(1)
		}
		outputFile = defaultPath
		fmt.Println("No output specified. Using:", outputFile)
	}
	if filepath.Ext(outputFile) != ".jsonl" {
		fmt.Println(red + "ERROR: Output file must have a '.jsonl' extension." + reset)
		os.Exit(1)
	}

	endTime := time.Now().Add(duration)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	//we need two cpu readings to caculate cpu usage
	prevCpuSnap, err := metrics.ReadCPUStatsRaw()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(green + "SNAP RUN STARTED" + reset)
	for now := range ticker.C {
		if now.After(endTime) {
			fmt.Println(green + "SNAP RUN COMPLETED" + reset)
			break
		}

		hasError := output.TakeSnapshot(&prevCpuSnap, outputFile, now)
		if hasError {
			fmt.Println(red + "ERROR: Snapshot failed - aborting." + reset)
			break
		}
		fmt.Println(blue + "SNAP: " + now.Local().Format("2006-01-02 15:04:05") + reset)

	}
}

func getDefaultOutputPath() (string, error) {
	folder := "snapsys_runs"
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := "snaprun_" + timestamp + ".jsonl"
	return filepath.Join(folder, filename), nil
}
