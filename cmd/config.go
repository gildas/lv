package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/kubectl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Initialize configure the logger and load the Viper Configuration
func Initialize(cmd *cobra.Command) (err error) {
	InitializeLogger(cmd)
	return InitializeConfiguration(cmd)
}

// InitializeLogger configures the logger based on the command line flags and environment variables
func InitializeLogger(cmd *cobra.Command) {
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

// InitializeConfiguration loads the configuration file and profiles
func InitializeConfiguration(cmd *cobra.Command) (err error) {
	viper.SetConfigType("yaml")
	if cmd != nil && cmd.Root().PersistentFlags().Changed("config") {
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

	_ = viper.BindPFlag("color", RootCmd.PersistentFlags().Lookup("color"))
	_ = viper.BindPFlag("follow", RootCmd.PersistentFlags().Lookup("follow"))
	_ = viper.BindPFlag("obfuscationKey", RootCmd.PersistentFlags().Lookup("key"))
	_ = viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("timezone", RootCmd.PersistentFlags().Lookup("time"))
	viper.SetDefault("color", true)
	viper.SetDefault("follow", false)
	viper.SetDefault("output", "long")
	viper.SetDefault("timezone", "local")

	viper.SetEnvPrefix("LV")
	viper.AutomaticEnv() // read in environment variables that match

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errors.Join(errors.New("Failed to read config file"), err)
		}
	}
	if err := kubectl.InitializeSelectors(cmd); err != nil {
		return errors.Join(errors.New("Failed to initialize selectors"), err)
	}
	return nil
}
