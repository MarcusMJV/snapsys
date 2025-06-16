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

const (
	green = "\033[32m"
	blue  = "\033[34m"
	red   = "\033[31m"
	reset = "\033[0m"
)

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
