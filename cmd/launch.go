package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// launchCmd represents a command to launch the updater service, periodically downloading CSV files
// from a website and updating a MySQL database with those values.
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch the updater service",
	Long: `Launch the updater service which periodically downloads CSV files from a website,
and updates a MySQL database with those values. The service uses a blue/green
deployment strategy using alternating tables to update the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		bindAllFlags(cmd)

		if err := viper.Unmarshal(&config); err != nil {
			logger.Error("unable to decode into struct", "error", err)
			os.Exit(1)
		}

		makeLogger(config.Logger.Level, config.Logger.Format)

		logger.Info("starting updater service", "config", config)
	},
}
