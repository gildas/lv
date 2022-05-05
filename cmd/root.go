package cmd

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Options caries options for this application
type Options struct {
	ConfigFile     string
	LogDestination string
	Verbose        bool
	Debug          bool
}

// CmdOptions contains the global options
var CmdOptions Options

// Log is the main logger
var Log *logger.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     APP,
	Short:   "Log viewer for Bunyan format",
	Version: Version(),
}

func init() {
	cobra.OnInitialize(initConfig)

	if home := core.GetEnvAsString("XDG_CONFIG_HOME", ""); len(home) > 0 {
		rootCmd.PersistentFlags().StringVarP(&CmdOptions.ConfigFile, "config", "c", "", fmt.Sprintf("config file (default is %s/um-cli/config)", home))
	} else {
		rootCmd.PersistentFlags().StringVarP(&CmdOptions.ConfigFile, "config", "c", "", "config file (default is $HOME/.um-cli)")
	}
	rootCmd.PersistentFlags().StringVar(&CmdOptions.LogDestination, "log", "", "where logs are writen if given (by default, no log is generated)")
	rootCmd.PersistentFlags().BoolVar(&CmdOptions.Debug, "debug", false, "forces logging at DEBUG level")
	rootCmd.PersistentFlags().BoolVarP(&CmdOptions.Verbose, "verbose", "v", false, "runs verbosely if set")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
	if Log != nil {
		Log.Flush()
	}
}

// initConfig reads config files and environment variable
func initConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	_ = viper.BindPFlag("DEBUG", rootCmd.Flags().Lookup("debug"))
	_ = viper.BindPFlag("VERBOSE", rootCmd.Flags().Lookup("verbose"))
	_ = viper.BindPFlag("LOG_DESTINATION", rootCmd.Flags().Lookup("log"))

	if len(CmdOptions.ConfigFile) > 0 { // Use config file from the flag.
		viper.SetConfigFile(CmdOptions.ConfigFile)
	} else if xdg := core.GetEnvAsString("XDG_CONFIG_HOME", ""); len(xdg) > 0 {
		viper.AddConfigPath(path.Join(xdg, APP))
		viper.SetConfigName("config")
	} else {
		viper.AddConfigPath(path.Join(home, ".config", APP))
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".um-cli")
		err = viper.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigName(".gum-cli")
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				viper.SetConfigType("dotenv")
				viper.SetConfigName(".env")
				err = viper.ReadInConfig()
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					err = nil // It's ok if the config file does not exist
				}
			}
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	os.Setenv("DEBUG", strconv.FormatBool(viper.GetBool("DEBUG")))
	os.Setenv("LOG_DESTINATION", viper.GetString("LOG_DESTINATION"))
	if len(viper.GetString("LOG_DESTINATION")) == 0 {
		if viper.GetBool("DEBUG") {
			Log = logger.Create(APP, APP+".log")
		} else {
			Log = logger.Create(APP, &logger.NilStream{})
		}
	} else {
		Log = logger.Create(APP)
	}

	Log.Infof(strings.Repeat("-", 80))
	Log.Infof("Starting %s version %s", APP, Version())
	Log.Infof("Configuration file: %s", viper.ConfigFileUsed())
	Log.Infof("Log Destination: %s", Log)
	Log.Infof("Verbose: %t", viper.GetBool("VERBOSE"))
}
