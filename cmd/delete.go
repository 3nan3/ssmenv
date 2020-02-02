package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/env"
	"github.com/3nan3/ssmenv/paramstore"
	"github.com/3nan3/ssmenv/util"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete environment variables",
	Args: func(cmd *cobra.Command, args []string) error {
		if deleteDryrun {
			deleteDiff = "all"
		} else if !util.SliceContains(deleteDiffs, deleteDiff) {
			return fmt.Errorf("hoge")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client := paramstore.New(getPath(), getEmptyPattern())

		var oldenvs *env.Env
		if deleteDiff != "no" {
			var err error
			oldenvs, err = client.GetEnvs()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}

		if !deleteDryrun {
			var err error
			deleteEnvVars, err = client.DeleteEnvs(deleteEnvVars)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}

		if deleteDiff != "no" {
			envs := env.New()
			for _, name := range deleteEnvVars {
				envs.PutEnv(name, "")
			}
			env.PrintDiff(oldenvs, envs, deleteDiff)
		}
	},
}

var (
	deleteEnvVars []string
	deleteDryrun bool
	deleteDiff string

	deleteDiffs = []string{"no", "all", "key"}
)

func init() {
	deleteCmd.Flags().StringSliceVarP(&deleteEnvVars, "env", "e", []string{}, "Environment variable name to delete (e.g. '-e ENV_VAR')")
	deleteCmd.Flags().BoolVar(&deleteDryrun, "dry-run", false, "Show differences and do nothing")
	deleteCmd.Flags().StringVar(&deleteDiff, "diff", "no", "Show delete result (e.g. \"--diff\", \"--diff=key\" )")
	deleteCmd.Flags().Lookup("diff").NoOptDefVal = "all"

	rootCmd.AddCommand(deleteCmd)
}
