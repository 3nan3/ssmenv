package cmd

import (
	"os"
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/env"
	"github.com/3nan3/ssmenv/paramstore"
	"github.com/3nan3/ssmenv/util"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Put environment variables",
	Args: func(cmd *cobra.Command, args []string) error {
		if putDryrun {
			putDiff = "all"
		} else if !util.SliceContains(putDiffs, putDiff) {
			return fmt.Errorf("\"diff\" must be one of the following: %s", deleteDiffs)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {		
		// load dotenv file
		envs := env.New()
		for _, file := range putEnvFiles {
			err := envs.LoadDotenv(file)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}

		// load env variables
		err := envs.LoadEnvVars(strings.Join(putEnvVars, "\n"))
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		// put env variables
		client := paramstore.New(getPath(), getEmptyPattern())
		var oldenvs *env.Env
		if !putDryrun {
			oldenvs, err = client.PutEnvs(envs)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		} else {
			oldenvs, err = client.GetEnvs()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}
		
		if putDiff != "no" {
			env.PrintDiff(oldenvs, envs, putDiff)
		}
	},
}
var (
	putEnvVars []string
	putEnvFiles []string
	putDryrun bool
	putDiff string

	putDiffs = []string{"no", "all", "key"}
)

func init() {
	putCmd.Flags().StringSliceVarP(&putEnvVars, "env", "e", []string{}, "a environment variable to put (e.g. '-e ENV_VAR=value')")
	putCmd.Flags().StringSliceVarP(&putEnvFiles, "env-file", "f", []string{}, "dotenv file to put")
	putCmd.Flags().BoolVar(&putDryrun, "dry-run", false, "show differences and do nothing")
	desc := fmt.Sprintf("show delete result (Select format: %s)", deleteDiffs)
	putCmd.Flags().StringVar(&putDiff, "diff", "no", desc)
	putCmd.Flags().Lookup("diff").NoOptDefVal = "all"

	rootCmd.AddCommand(putCmd)
}
