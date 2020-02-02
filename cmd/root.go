package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use:   "ssmenv",
	Short: "Use keys and values of AWS SSM Parameter Store as environment variables",
}
var path string
var emptyPattern string


func Execute() {
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
