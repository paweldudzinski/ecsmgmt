package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists available resources.",
	Long:  `Lists available resources: clusters, services or tasks.`,
	// No run command will print the help message
}

func init() {
	rootCmd.AddCommand(listCmd)
}
