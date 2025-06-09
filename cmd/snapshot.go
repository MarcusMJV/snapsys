/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/MarcusMJV/snapsys.git/internal/metrics"
	"github.com/MarcusMJV/snapsys.git/internal/output"
	"github.com/spf13/cobra"
)

type Snapshot struct {
}

var interval time.Duration
var duration time.Duration
var outputFile string

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snapshot called")
		runSnapshot()
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().DurationVar(&interval, "interval", 5*time.Second, "Interval between snapshots")
	snapshotCmd.Flags().DurationVar(&duration, "duration", 1*time.Minute, "Total duration to run snapshots")
	snapshotCmd.Flags().StringVar(&outputFile, "output", "snapshot.json", "Output file path")
}

func runSnapshot() {

	endTime := time.Now().Add(duration)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	prevCpuSnap, err := metrics.ReadCPUStats()
	if err != nil {
		fmt.Println(err)
	}

	var snapshots []output.Snapshot

	for now := range ticker.C {
		if now.After(endTime) {
			fmt.Println("Sbnap run completed")
			break
		}

		cpuSnap, err := metrics.ReadCPUStats()
		if err != nil {
			fmt.Println(err)
		}

		cpuUsage := metrics.CalculateCpuUsage(prevCpuSnap, cpuSnap)
		prevCpuSnap = cpuSnap
		fmt.Println("CPU Percentage: ", cpuUsage)

		snapshots = append(snapshots, output.Snapshot{Timestamp: now, CPUUsage: cpuUsage})

	}

	if err := output.WriteSnapshotsToFile(snapshots, outputFile); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Output written to: ", outputFile)
	}

}
