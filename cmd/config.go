package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Initialize configure the logger and load the Viper Configuration
func Initialize(cmd *cobra.Command) (err error) {
	initializeLogger(cmd)
	return initializeConfiguration(cmd)
}

// initializeLogger configures the logger based on the command line flags and environment variables
func initializeLogger(cmd *cobra.Command) {
	log := logger.Must(logger.FromContext(cmd.Context()))

	// Use persistent flags instead
	if cmd.Root().PersistentFlags().Changed("log") {
		log.ResetDestinations(cmd.Root().PersistentFlags().Lookup("log").Value.String())
	}

	log.Infof("%s", strings.Repeat("-", 80))
	log.Infof("Starting %s v%s (%s)", cmd.Root().Name(), cmd.Root().Version, runtime.GOARCH)
	log.Infof("Log Destination: %s", log)
	if cmd.Root().PersistentFlags().Changed("debug") && cmd.Root().PersistentFlags().Lookup("debug").Value.String() == "true" {
		log.SetFilterLevel(logger.DEBUG)
		log.Infof("Debug was turned on by the --debug flag")
	}
}

// initializeConfiguration loads the configuration file and profiles
func initializeConfiguration(cmd *cobra.Command) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context()))

	viper.SetConfigType("yaml")
	if cmd.Root().PersistentFlags().Changed("config") {
		viper.SetConfigFile(cmd.Root().PersistentFlags().Lookup("config").Value.String())
	} else if configDir, _ := os.UserConfigDir(); len(configDir) > 0 {
		viper.AddConfigPath(filepath.Join(configDir, "logviewer"))
		viper.SetConfigName("config.yaml")
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".logviewer")
	}

	_ = viper.BindPFlag("local", RootCmd.PersistentFlags().Lookup("local"))
	_ = viper.BindPFlag("timezone", RootCmd.PersistentFlags().Lookup("time"))
	_ = viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("color", RootCmd.PersistentFlags().Lookup("color"))
	_ = viper.BindPFlag("obfuscationKey", RootCmd.PersistentFlags().Lookup("key"))
	viper.SetDefault("local", false)
	viper.SetDefault("timezone", "UTC")
	viper.SetDefault("output", "long")
	viper.SetDefault("color", true)

	viper.SetEnvPrefix("LV")
	_ = viper.BindEnv("local", "obfuscationKey")
	viper.AutomaticEnv() // read in environment variables that match

	err = viper.ReadInConfig()
	if verr, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warnf("Config file not found: %s", verr)
	} else if err != nil {
		return errors.Join(errors.New("Failed to read config file"), err)
	} else {
		log.Infof("Config File: %s", viper.ConfigFileUsed())
	}
	return nil
}
