// Package config provides configuration management for the application.
package config

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

// MakeLogger creates a new slog logger based on the set configuration.
func (c *Config) MakeLogger() (*slog.Logger, error) {
	var slogLevel slog.Level
	switch strings.ToLower(c.Logger.Level) {
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
	switch strings.ToLower(c.Logger.Format) {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	default:
		return nil, errors.New("invalid log format, must be 'text' or 'json'")
	}

	return slog.New(handler), nil
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

// InitConfig initializes the Viper configuration by reading from a config file
// or environment variables.
func InitConfig(cfgPath string, logger *slog.Logger) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
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

// BindAllFlags binds all user-changed flags in a Cobra FlagSet to Viper configuration keys
func BindAllFlags(cmd *cobra.Command) {
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
