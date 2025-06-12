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
	snapshotCmd.Flags().StringVar(&outputFile, "output", "snapshot.jsonl", "Output file path")
}

func runSnapshot() {

	endTime := time.Now().Add(duration)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	prevCpuSnap, err := metrics.ReadCPUStatsRaw()
	if err != nil {
		fmt.Println(err)
	}

	for now := range ticker.C {
		if now.After(endTime) {
			fmt.Println("Snap run completed")
			break
		}

		cpuRaw, err := metrics.ReadCPUStatsRaw()
		if err != nil {
			fmt.Println(err)
		}

		cpuUsage := metrics.CalculateCpuUsage(prevCpuSnap, cpuRaw)
		cpuStats := metrics.CPUStats{UsagePct: cpuUsage, Raw: cpuRaw}
		prevCpuSnap = cpuRaw
		fmt.Println("CPU Percentage: ", cpuStats.UsagePct)

		memoryStats, err := metrics.ReadMemStats()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Memory Percentage: ", memoryStats.UsagePct)

		diskStas, err := metrics.GetAllDisks()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Got Disk Stats")

		snapshot := output.Snapshot{
			Timestamp: now,
			CPU:       cpuStats,
			Memory:    memoryStats,
			Disks:     diskStas,
		}
		err = snapshot.AppendSnapshotJSONL(outputFile)

		if err != nil {
			fmt.Println(err)
		}

	}

	fmt.Println("file created")

}
