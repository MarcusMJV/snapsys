/*
SnapSys - Lightweight System Benchmarking Tool
Copyright © 2025 Marcus Vorster
Released under the MIT License
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:   "snapsys",
	Short: "SnapSys is a lightweight system benchmarking tool",
	Long: `SnapSys is a terminal-based benchmarking utility that captures snapshots of 
CPU, memory, and disk usage over time. It's ideal for developers, sysadmins, 
and power users who need a fast, minimal way to log system performance.

Usage Examples:

  snapsys snapshot --duration 60s --interval 5s --output stats.jsonl

This command captures system metrics every 5 seconds for 60 seconds and
writes the output to 'stats.jsonl' in JSON Lines format.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("SnapSys version: {{.Version}}\n")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
