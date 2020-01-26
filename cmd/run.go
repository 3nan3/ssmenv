package cmd

import (
	"os"
	"os/exec"
	"syscall"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/paramstore"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := paramstore.New(getPath())
		envs, err := client.GetEnvs()
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		envs.ApplyEnv()
		path, err := exec.LookPath(args[0])
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		err = syscall.Exec(path, args, os.Environ())
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}		
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
