// Package cmd provides the launch point for the updater service, using Cobra for CLI and Viper for
// configuration management to set up the service's parameters and run it.
package cmd

import (
	"log/slog"
	"os"

	cfg "github.com/lorendsnow/updater/internal/config"
	"github.com/spf13/cobra"
)

/*
 *==================================================================================================
 * Package Constants
 *==================================================================================================
 */

var (
	cfgFile string
	config  cfg.Config
	logger  = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	rootCmd = &cobra.Command{
		Use:   "updater",
		Short: "A database updater service",
		Long: `Updater is a service that periodically downloads CSV files from a website, and
		updates a MySQL database with those values.`,
	}
)

/*
 *==================================================================================================
 * Public Functions
 *==================================================================================================
 */

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

/*
 *==================================================================================================
 * Private Functions
 *==================================================================================================
 */

// init sets up the Cobra CLI interface, and a Viper configuration
func init() {
	cobra.OnInitialize(initViper)

	rootCmd.AddCommand(launchCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "path to config file")
	rootCmd.PersistentFlags().String("host", "", "MySQL host")
	rootCmd.PersistentFlags().Int("port", 0, "MySQL port")
	rootCmd.PersistentFlags().String("user", "", "MySQL user")
	rootCmd.PersistentFlags().String("pass", "", "MySQL password")
	rootCmd.PersistentFlags().String("name", "", "MySQL database name")
	rootCmd.PersistentFlags().String("interval", "", "check interval")
	rootCmd.PersistentFlags().StringArray("csv", []string{}, "CSV URLs")
	rootCmd.PersistentFlags().String("blue-table", "", "blue table name")
	rootCmd.PersistentFlags().String("green-table", "", "green table name")
	rootCmd.PersistentFlags().String("timeout", "", "HTTP timeout")
	rootCmd.PersistentFlags().Int("retries", 0, "HTTP retries")
	rootCmd.PersistentFlags().String(
		"log-level",
		"",
		"log level (one of debug, info, warn or error)",
	)
	rootCmd.PersistentFlags().String(
		"log-format",
		"",
		"log format (one of json or text)",
	)
}

// initViper runs the Viper initialization function from the config package.
func initViper() {
	cfg.InitConfig(cfgFile, logger)
}
