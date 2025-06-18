package cmd

import (
	"os"

	cfg "github.com/lorendsnow/updater/internal/config"
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
		cfg.BindAllFlags(cmd)

		if err := viper.Unmarshal(&config); err != nil {
			logger.Error("unable to decode into struct", "error", err)
			os.Exit(1)
		}

		appLogger, err := config.MakeLogger()
		if err != nil {
			config.Logger.Level = "info"
			config.Logger.Format = "text"
			logger.Error(
				"unable to create application logger, using default logging configuration",
				"error",
				err,
			)
		}

		if appLogger != nil {
			logger = appLogger
		}

		logger.Info("starting updater service", "config", config)
	},
}
