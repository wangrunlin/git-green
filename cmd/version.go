package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Green",
	RunE: func(cmd *cobra.Command, args []string) error {
		printGitGreenVersion()
		return nil
	},
	DisableFlagParsing: true,
}

func printGitGreenVersion() {
	fmt.Println("git-green version v0.0.2")
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
