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
		if dryrun {
			diff = "all"
		} else if !util.SliceContains(diffs, diff) {
			return fmt.Errorf("hoge")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {		
		// load dotenv file
		envs := env.New()
		for _, file := range envFiles {
			err := envs.LoadDotenv(file)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}

		// load env variables
		err := envs.LoadEnvVars(strings.Join(envVars, "\n"))
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		// put env variables
		client := paramstore.New(getPath(), getEmptyPattern())
		var oldenvs *env.Env
		if !dryrun {
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
		
		if diff != "no" {
			env.PrintDiff(oldenvs, envs, diff)			
		}
	},
}
var (
	envVars []string
	envFiles []string
	dryrun bool
	diff string

	diffs = []string{"no", "all", "key"}
)

func init() {
	putCmd.Flags().StringSliceVarP(&envVars, "env", "e", []string{}, "environment variable (e.g. '-e ENV_VAR=value')")
	putCmd.Flags().StringSliceVarP(&envFiles, "env-file", "f", []string{}, "dotenv file")
	putCmd.Flags().BoolVar(&dryrun, "dry-run", false, "")
	putCmd.Flags().StringVar(&diff, "diff", "no", "")
	putCmd.Flags().Lookup("diff").NoOptDefVal = "all"

	rootCmd.AddCommand(putCmd)
}
