package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/paramstore"
)

var runCmd = &cobra.Command{
	Use:   "run <any command>",
	Short: "Execute command with environment variables applied",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires any command to execute")
		}
		return nil
	},
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "-h" || args[0] == "--help" {
			cmd.Help()
			os.Exit(0)
		}

		client := paramstore.New(getPath())
		envs, err := client.GetEnvs()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		envs.ApplyEnv()
		path, err := exec.LookPath(args[0])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		err = syscall.Exec(path, args, os.Environ())
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}		
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
