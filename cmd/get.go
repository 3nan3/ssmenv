package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/3nan3/ssmenv/paramstore"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := paramstore.New("/dotenv/development")			
		env, err := client.GetEnv(args[0])
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
		for name, value := range env {
			fmt.Printf("%s='%s'\n", name, value)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
