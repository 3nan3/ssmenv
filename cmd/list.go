package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/paramstore"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display list of environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		client := paramstore.New(getPath(), getEmptyPattern())
		envs, err := client.GetEnvs()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		err = envs.PrintAll(listExport); if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}
var listExport bool

func init() {
	listCmd.Flags().BoolVar(&listExport, "export", false, "print with export at the beginning")

	rootCmd.AddCommand(listCmd)
}
