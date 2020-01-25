package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use:   "ssmenv <command>",
	Short: "",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("run")
		os.Exit(0)
	},
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
