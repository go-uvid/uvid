package cmd

import (
	"luvsic3/uvid/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

const DSN = "uvid.db"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the web server (default port 3000), use PORT environment variable to override the port",
	Run: func(cmd *cobra.Command, args []string) {
		api.New(DSN).Start()
	},
}
