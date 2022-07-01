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
	Input          string
	LogLevel       string
	Filter         string
	LocalTime      bool
	UseColors      bool
	UsePager       bool
	Strict         bool
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
	Run:     runRootCommand,
}

func init() {
	cobra.OnInitialize(initConfig)

	if home := core.GetEnvAsString("XDG_CONFIG_HOME", ""); len(home) > 0 {
		rootCmd.PersistentFlags().StringVarP(&CmdOptions.ConfigFile, "config", "c", "", fmt.Sprintf("config file (default is %s/um-cli/config)", home))
	} else {
		rootCmd.PersistentFlags().StringVarP(&CmdOptions.ConfigFile, "config", "c", "", "config file (default is $HOME/.um-cli)")
	}
	rootCmd.PersistentFlags().StringVar(&CmdOptions.Input, "input", "i", "stdin", "read the log from the given file. If absent or file is \"-\", stdin is read.")
	rootCmd.PersistentFlags().StringVar(&CmdOptions.LogLevel, "level", "l", "INFO", "Only shows log entries with a level at or above the given value.")
	rootCmd.PersistentFlags().StringVar(&CmdOptions.Filter, "filter", "f", "", "Run each log message through the filter.")
	rootCmd.PersistentFlags().BoolVar(&CmdOptions.Strict, "strict", "", false, "Suppress all but legal Bunyan JSON log lines. By default non-JSON, and non-Bunyan lines are passed through.")
	rootCmd.PersistentFlags().BoolVar(&CmdOptions.UsePager, "pager", "", false, "Pipe output into `less` (or $PAGER if set), if stdout is a TTY. This overrides $BUNYAN_NO_PAGER.")
	rootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "color", "", true, "Colorize output. Defaults to try if output stream is a TTY.")
	rootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "no-color", "", false, "Force no coloring.")
	rootCmd.PersistentFlags().StringVar(&CmdOptions.Output, "output", "o", "long", "")
	rootCmd.PersistentFlags().StringVar(&CmdOptions.Output, "output", "o", "long", "")

	// LogLevel should also support: https://github.com/gildas/go-logger#setting-the-filterlevel

	// --strict suppresses all but legal Bunyan log entries. By default, non-Bunyan entries are passed through.
	/*
		   p('  -c, --condition CONDITION');
		   p('                Run each log message through the condition and');
		   p('                only show those that return truish. E.g.:');
		   p('                    -c \'this.pid == 123\'');
		   p('                    -c \'this.level == DEBUG\'');
		   p('                    -c \'this.msg.indexOf("boom") != -1\'');
		   p('                "CONDITION" must be legal JS code. `this` holds');
		   p('                the log record. The TRACE, DEBUG, ... FATAL values');
		   p('                are defined to help with comparing `this.level`.');
		   How about some Go Template?
		   p('Output options:');
		   p('  --no-pager    Do not pipe output into a pager.');
		   p('  --no-color    Force no coloring (e.g. terminal doesn\'t support it)');
		   p('  -o, --output MODE');
		   p('                Specify an output mode/format. One of');
		   p('                  long: (the default) pretty');
		   p('                  json: JSON output, 2-space indent');
		   p('                  json-N: JSON output, N-space indent, e.g. "json-4"');
		   p('                  bunyan: 0 indented JSON, bunyan\'s native format');
		   p('                  inspect: node.js `util.inspect` output');
		   p('                  short: like "long", but more concise');
		   p('                  simple: level, followed by "-" and then the message');
		                        html: generate an html page
								serve, server: starts a web server to give the html page. Should be dynamic, etc.
		   p('  -j            shortcut for `-o json`');
		   p('  -0            shortcut for `-o bunyan`');
		   p('  -L, --time local');
		   p('                Display time field in local time, rather than UTC.');
		   p('');
		   p('Environment Variables:');
		   p('  BUNYAN_NO_COLOR    Set to a non-empty value to force no output ');
		   p('                     coloring. See "--no-color".');
		   p('  BUNYAN_NO_PAGER    Disable piping output to a pager. ');
		   p('                     See "--no-pager".');
	*/
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

// runRootCommand executes the Root Command
func runRootCommand(cmd *cobra.Command, args []string) {
	Log.Record("options", CmdOptions).Debugf("Option Flags")
}
