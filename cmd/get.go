package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/paramstore"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := paramstore.New(getPath())
		envs, err := client.GetEnv(args[0])
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		err = envs.Stdout(); if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
