package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/crissyfield/troll-a/cmd"
)

// Version will be set during build.
var Version = "(unknown)"

// CmdRoot defines the CLI root command.
var CmdRoot = &cobra.Command{
	Use:               "troll-a",
	Long:              "...",
	Args:              cobra.NoArgs,
	Version:           Version,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: setup,
}

// Initialize CLI options.
func init() {
	// Logging
	CmdRoot.PersistentFlags().String("log-level", "info", "verbosity of logging output")
	CmdRoot.PersistentFlags().Bool("log-as-json", false, "change logging format to JSON")

	// Subcommands
	CmdRoot.AddCommand(cmd.CmdTest)
}

// setup will set up configuration management and logging.
//
// Configuration options can be set via the command line, via a configuration file (in the current folder, at
// "/etc/troll-a/config.yaml" or at "~/.config/troll-a/config.yaml"), and via environment variables (all
// uppercase and prefixed with "TROLLA_").
func setup(cmd *cobra.Command, _ []string) error {
	// Connect all options to Viper
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return fmt.Errorf("bind command line flags: %w", err)
	}

	// Environment variables
	viper.SetEnvPrefix("TROLLA")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// Configuration file
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/troll-a")
	viper.AddConfigPath(os.Getenv("HOME") + "/.config/troll-a")
	viper.AddConfigPath(".")

	// Configuration file
	if err := viper.ReadInConfig(); err != nil {
		// Don't fail if config not found
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("read config file: %w", err)
		}
	}

	// Logging
	var level slog.Level

	err = level.UnmarshalText([]byte(viper.GetString("log-level")))
	if err != nil {
		return fmt.Errorf("parse log level: %w", err)
	}

	var handler slog.Handler

	if viper.GetBool("log-as-json") {
		// Use JSON handler
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	} else {
		// Use text handler
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}

	slog.SetDefault(slog.New(handler))

	return nil
}

// main is the main entry point of the command.
func main() {
	if err := CmdRoot.Execute(); err != nil {
		slog.Error("Unable to execute command", slog.Any("error", err))
	}
}
