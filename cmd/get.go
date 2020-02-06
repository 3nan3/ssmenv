package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/paramstore"
)

var getCmd = &cobra.Command{
	Use:   "get <env_name>",
	Short: "Display single environment variable",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires environment variable name")
		}
		return nil
	},
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := paramstore.New(getPath(), getEmptyPattern())
		envs, err := client.GetEnv(args[0])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		err = envs.PrintAll(); if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
