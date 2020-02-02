package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

const version = "v0.1.0"

var rootCmd = &cobra.Command {
	Use:   "ssmenv",
	Short: "Use keys and values of AWS SSM Parameter Store as environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		if rootVersion {
			printVersion()
		} else {
			cmd.Help()
		}
	},
}
var (
	path string
	emptyPattern string

	rootVersion bool
)

func Execute() {
	rootCmd.Flags().BoolVarP(&rootVersion, "version", "v", false, "show ssmenv version")

	rootCmd.Execute()
}

func getPath() string {
	if path == "" {
		path = os.Getenv("SSMENV_PATH")
	}
	return path
}

func getEmptyPattern() string {
	if emptyPattern == "" {
		emptyPattern = os.Getenv("SSMENV_EMPTY_PATTERN")
		if emptyPattern == "" {
			emptyPattern = "🈳"
		}
	}
	return emptyPattern
}

func printVersion() {
	fmt.Printf("ssmenv %s\n", version)
}
