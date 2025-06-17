// Package cmd provides the launch point for the updater service, using Cobra for CLI and Viper for
// configuration management to set up the service's parameters and run it.
package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
 *==================================================================================================
 * Package Constants
 *==================================================================================================
 */

var (
	cfgFile string
	config  Config
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
 * Config Struct
 *==================================================================================================
 */

// Config holds configuration values for the updater service.
type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`

	Service struct {
		CheckInterval string   `mapstructure:"check-interval"`
		CSVUrls       []string `mapstructure:"csv-urls"`
		BlueTable     string   `mapstructure:"blue-table"`
		GreenTable    string   `mapstructure:"green-table"`
	} `mapstructure:"service"`

	HTTP struct {
		Timeout string `mapstructure:"timeout"`
		Retries int    `mapstructure:"retries"`
	} `mapstructure:"http"`

	Logger struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logger"`
}

/*
 *==================================================================================================
 * FlagName Enum
 *==================================================================================================
 */

// FlagName represents the Cobra flag names used in the application.
type FlagName int

const (
	Host FlagName = iota
	Port
	User
	Pass
	Name
	Interval
	CSV
	BlueTable
	GreenTable
	Timeout
	Retries
	LogLevel
	LogFormat
)

// String returns the string representation of the FlagName.
func (f FlagName) String() string {
	switch f {
	case Host:
		return "host"
	case Port:
		return "port"
	case User:
		return "user"
	case Pass:
		return "pass"
	case Name:
		return "name"
	case Interval:
		return "interval"
	case CSV:
		return "csv"
	case BlueTable:
		return "blue-table"
	case GreenTable:
		return "green-table"
	case Timeout:
		return "timeout"
	case Retries:
		return "retries"
	case LogLevel:
		return "log-level"
	case LogFormat:
		return "log-format"
	default:
		return ""
	}
}

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
	cobra.OnInitialize(initConfig)

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

// initConfig initializes the Viper configuration by reading from a config file
// or environment variables.
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	viper.SetEnvPrefix("UPDATER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Error("error reading config file", "error", err)
			os.Exit(1)
		}
	}
}

// bindAllFlags binds all user-changed flags in a Cobra FlagSet to Viper configuration keys
func bindAllFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		var viperName string

		switch flag.Name {
		case Host.String():
			viperName = "database.host"
		case Port.String():
			viperName = "database.port"
		case User.String():
			viperName = "database.username"
		case Pass.String():
			viperName = "database.password"
		case Name.String():
			viperName = "database.name"
		case Interval.String():
			viperName = "service.check-interval"
		case CSV.String():
			viperName = "service.csv-urls"
		case BlueTable.String():
			viperName = "service.blue-table"
		case GreenTable.String():
			viperName = "service.green-table"
		case Timeout.String():
			viperName = "http.timeout"
		case Retries.String():
			viperName = "http.retries"
		case LogLevel.String():
			viperName = "logger.level"
		case LogFormat.String():
			viperName = "logger.format"
		default:
			return
		}

		if viperName != "" {
			viper.BindPFlag(viperName, flag)
		}
	})
}

// makeLogger creates a new slog logger based on the provided configuration.
func makeLogger(level string, output string) {
	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	var handler slog.Handler
	switch strings.ToLower(output) {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	default:
		slog.SetDefault(logger)
		return
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}
