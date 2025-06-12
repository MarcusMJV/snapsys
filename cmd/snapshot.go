/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sync"
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

		var wg sync.WaitGroup
		wg.Add(3)
		errChan := make(chan error, 1)
		fmt.Println("channels made")
		cpuChan := make(chan metrics.CPUStats, 1)
		memChan := make(chan metrics.MemoryStats, 1)
		diskChan := make(chan metrics.DiskMap, 1)

		go func() {
			defer wg.Done()
			cpu, err := metrics.ReadCPU(&prevCpuSnap)
			fmt.Println("cpu read")
			if err != nil {
				fmt.Println(err)
				errChan <- err
			}
			cpuChan <- cpu
		}()
		go func() {
			defer wg.Done()
			memory, err := metrics.ReadMemStats()
			fmt.Println("memory read")
			if err != nil {
				errChan <- err
			}
			memChan <- memory
		}()
		go func() {
			defer wg.Done()
			disks, err := metrics.GetAllDisks()
			fmt.Println("disks read")
			if err != nil {
				errChan <- err
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
			fmt.Println("Aborting snapshot run.")
			break
		}

		snapshot := output.Snapshot{
			Timestamp: now,
			CPU:       <-cpuChan,
			Memory:    <-memChan,
			Disks:     <-diskChan,
		}

		err = snapshot.AppendSnapshotJSONL(outputFile)
		fmt.Println("wrote snapshot")
		if err != nil {
			fmt.Println(err)
		}

	}
}
