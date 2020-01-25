package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/paramstore"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		client := paramstore.New("/dotenv/development")			
		envs, err := client.GetEnvs()
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		for name, value := range envs {
			fmt.Printf("%s='%s'\n", name, value)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
