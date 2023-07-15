package cmd

import (
	"github.com/rick-you/uvid/api"

	"github.com/spf13/cobra"
)

const DSN = "uvid.db"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "uvid",
	Long: "Start uvid server (default port 3000), use PORT environment variable to override the port",
	Run: func(cmd *cobra.Command, args []string) {
		api.New(DSN).Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
