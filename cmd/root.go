package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type OutputOptions struct {
	LogLevel  string
	Filter    string
	Output    *flags.EnumFlag
	Location  *time.Location
	UseColors bool
}

// CmdOptions contains the global options
var CmdOptions struct {
	OutputOptions
	Completion     *flags.EnumFlag
	ConfigFile     string
	LogDestination string
	Timezone       string
	UsePager       bool
	Verbose        bool
	Debug          bool
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Short: "pretty-print logviewer logs from stdin or file(s)",
	Long:  "logviewer is a simple and fast JSON log viewer. It reads log entries from given files or stdin and pretty-prints them to stdout.",
	RunE:  runRootCommand,
}

// Execute run the command
func Execute(context context.Context) error {
	return RootCmd.ExecuteContext(context)
}

func init() {
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	CmdOptions.Output = flags.NewEnumFlag("+long", "logviewer", "short", "simple", "html", "serve", "server")
	CmdOptions.Completion = flags.NewEnumFlag("bash", "zsh", "fish", "powershell", "help")
	RootCmd.PersistentFlags().Var(CmdOptions.Completion, "completion", "Generates completion script for bash, zsh, fish, or powershell")
	RootCmd.PersistentFlags().StringVar(&CmdOptions.ConfigFile, "config", "", fmt.Sprintf("config file (default is %s)", filepath.Join(configDir, "logviewer", "config.yaml")))
	RootCmd.PersistentFlags().StringVar(&CmdOptions.LogLevel, "level", "", "Only shows log entries with a level at or above the given value.")
	RootCmd.PersistentFlags().StringVarP(&CmdOptions.Filter, "filter", "f", "", "Run each log message through the filter.")
	RootCmd.PersistentFlags().StringVarP(&CmdOptions.Filter, "condition", "c", "", "Run each log message through the filter.")
	RootCmd.PersistentFlags().BoolP("local", "L", false, "Display time field in local time, rather than UTC.")
	RootCmd.PersistentFlags().StringVar(&CmdOptions.Timezone, "time", "", "Display time field in the given timezone.")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UsePager, "no-pager", true, "Do not pipe output into a pager. By default, the output is piped throug `less` (or $PAGER if set), if stdout is a TTY")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "no-color", false, "Do not colorize output. By default, the output is colorized if stdout is a TTY")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "color", true, "Colorize output always, even if the output stream is not a TTY.")
	RootCmd.PersistentFlags().VarP(CmdOptions.Output, "output", "o", "output mode/format. One of long, json, json-N, logviewer, inspect, short, simple, html, serve, server")
	RootCmd.PersistentFlags().StringVar(&CmdOptions.LogDestination, "log", "", "where logs are writen if given (by default, no log is generated)")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.Debug, "debug", false, "forces logging at DEBUG level")
	RootCmd.PersistentFlags().BoolVarP(&CmdOptions.Verbose, "verbose", "v", false, "runs verbosely if set")
	_ = RootCmd.RegisterFlagCompletionFunc(CmdOptions.Output.CompletionFunc("output"))
	_ = RootCmd.RegisterFlagCompletionFunc(CmdOptions.Completion.CompletionFunc("completion"))

	RootCmd.SilenceUsage = true
	_ = viper.BindPFlag("local", RootCmd.PersistentFlags().Lookup("local"))
	_ = viper.BindPFlag("timezone", RootCmd.PersistentFlags().Lookup("time"))
	_ = viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("color", RootCmd.PersistentFlags().Lookup("color"))
	viper.SetDefault("local", false)
	viper.SetDefault("timezone", "UTC")
	viper.SetDefault("output", "long")
	viper.SetDefault("color", true)

	cobra.OnInitialize(initConfig)
}

// initConfig reads config files and environment variable
func initConfig() {
	log := logger.Must(logger.FromContext(RootCmd.Context()))

	if len(CmdOptions.LogDestination) > 0 {
		log.ResetDestinations(CmdOptions.LogDestination)
	}

	log.Infof("%s", strings.Repeat("-", 80))
	log.Infof("Starting %s v%s (%s)", RootCmd.Name(), RootCmd.Version, runtime.GOARCH)
	log.Infof("Log Destination: %s", log)

	if CmdOptions.Debug {
		log.SetFilterLevel(logger.DEBUG)
		log.Infof("Debug was turned on by the --debug flag")
	}

	viper.SetConfigType("yaml")
	if len(CmdOptions.ConfigFile) > 0 { // Use config file from the flag.
		viper.SetConfigFile(CmdOptions.ConfigFile)
	} else if configDir, _ := os.UserConfigDir(); len(configDir) > 0 {
		viper.AddConfigPath(filepath.Join(configDir, "logviewer"))
		viper.SetConfigName("config.yaml")
	} else { // Old fashion way
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".logviewer")
	}

	viper.SetEnvPrefix("LV_")
	_ = viper.BindEnv("local")
	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()
	if verr, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warnf("Config file not found: %s", verr)
	} else if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
		fmt.Fprintf(os.Stderr, "Failed to read config file: %s\n", err)
		os.Exit(1)
	} else {
		log.Infof("Config File: %s", viper.ConfigFileUsed())
	}
}

// runRootCommand executes the Root Command
func runRootCommand(cmd *cobra.Command, args []string) (err error) {
	// Here we should read from stdin or from the files
	log := logger.Must(logger.FromContext(cmd.Context()))
	var reader *bufio.Reader

	if cmd.Flags().Changed("completion") {
		return generateCompletion(cmd, CmdOptions.Completion.Value)
	}

	CmdOptions.UseColors = isStdoutTTY() || viper.GetBool("color") || cmd.Flags().Changed("color")
	if cmd.Flags().Changed("no-color") {
		CmdOptions.UseColors = false
	}
	CmdOptions.UsePager = isStdoutTTY() && isStdinTTY()
	if cmd.Flags().Changed("no-pager") {
		CmdOptions.UsePager = false
	}
	CmdOptions.Output.Value = viper.GetString("output")

	if viper.GetBool("local") {
		log.Infof("Displaying local time")
		CmdOptions.Location = time.Local
	} else if CmdOptions.Location, err = ParseLocation(viper.GetString("timezone")); err != nil {
		log.Fatalf("Failed to load timezone %s: %s", viper.GetString("timezone"), err)
		return err
	}
	log.Infof("Displaying time at location: %s", CmdOptions.Location)

	if len(args) == 0 {
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatalf("Failed to open file %s: %s", args[0], err)
			return err
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	var outstream io.WriteCloser = os.Stdout

	if CmdOptions.UsePager {
		var closer func()

		if outstream, closer, err = GetPager(cmd.Context()); err != nil {
			log.Fatalf("Failed to get pager", err)
			return err
		}
		defer closer()
	}

	filters := MultiLogFilter{}

	if len(CmdOptions.LogLevel) > 0 {
		log.Infof("Adding log level filter at %s", CmdOptions.LogLevel)
		filters.Add(NewLevelLogFilter(CmdOptions.LogLevel))
	}
	if len(CmdOptions.Filter) > 0 {
		log.Infof("Adding filter: %s", CmdOptions.Filter)
		filter, err := NewConditionFilter(CmdOptions.Filter)
		if err != nil {
			log.Fatalf("Failed to create filter: %s", err)
			return err
		}
		filters.Add(filter)
		log.Infof("Added Filter: %#v", filter)
	}
	var filter LogFilter = filters.AsFilter()

	for {
		var line []byte

		line, err = ReadLine(reader)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Errorf("Failed to read line: %s", err)
			}
			break
		}

		if len(line) == 0 {
			continue
		}
		log.Infof("%s", string(line))
		var entry LogEntry

		if err := json.Unmarshal(line, &entry); err != nil {
			log.Errorf("Failed to parse JSON: %s", err)
			fmt.Fprintln(outstream, string(line))
			continue
		}
		if filter.Filter(cmd.Context(), entry) {
			output := strings.Builder{}

			entry.Write(cmd.Context(), &output, &CmdOptions.OutputOptions)
			if output.Len() > 0 {
				fmt.Fprintln(outstream, output.String())
			}
		}
	}
	if err != nil && !errors.Is(err, io.EOF) {
		log.Fatalf("Failed to read from input", err)
		return err
	}
	return nil
}
