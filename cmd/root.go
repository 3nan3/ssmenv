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


func Execute() {
	rootCmd.Execute()
}

func getPath() string {
	if path == "" {
		path = os.Getenv("SSMENV_PATH")
	}
	return path
}
