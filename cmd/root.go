package cmd

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "app"}

// Execute is the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(getStrategyCmd())
}
