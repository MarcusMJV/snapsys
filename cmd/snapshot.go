/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

type Snapshot struct {
}

var interval time.Duration
var duration time.Duration
var output string

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
	snapshotCmd.Flags().StringVar(&output, "output", "snapshot.json", "Output file path")
}

func runSnapshot() {
	fmt.Println("running snap shot")
	// fmt.Printf("running snapshot every %v for %v saving %v", interval, duration, output)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	done := time.After(duration)

	for {
		select {
		case <-done:
			fmt.Println("Snap run completed")
		case t := <-ticker.C:
			fmt.Println("Taking Snapshot: ", t)
		}

	}
}
