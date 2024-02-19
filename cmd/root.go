package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CmdOptions contains the global options
var CmdOptions struct {
	ConfigFile     string
	LogDestination string
	LogLevel       string
	Filter         string
	Output         string
	LocalTime      bool
	UseColors      bool
	UsePager       bool
	Strict         bool
	Verbose        bool
	Debug          bool
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Short: "pretty-print Bunyan logs from stdin or file(s)",
	Long:  "Bunyan is a simple and fast JSON log viewer. It reads log entries from given files or stdin and pretty-prints them to stdout.",
	RunE:  runRootCommand,
}

// Execute run the command
func Execute(context context.Context) error {
	return RootCmd.ExecuteContext(context)
}

func init() {
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	RootCmd.PersistentFlags().StringVarP(&CmdOptions.ConfigFile, "config", "c", "", fmt.Sprintf("config file (default is %s)", filepath.Join(configDir, "bunyan", "config.yaml")))
	RootCmd.PersistentFlags().StringVar(&CmdOptions.LogLevel, "level", "INFO", "Only shows log entries with a level at or above the given value.")
	RootCmd.PersistentFlags().StringVarP(&CmdOptions.Filter, "filter", "f", "", "Run each log message through the filter.")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.Strict, "strict", false, "Suppress all but legal Bunyan JSON log lines. By default non-JSON, and non-Bunyan lines are passed through.")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UsePager, "pager", false, "Pipe output into `less` (or $PAGER if set), if stdout is a TTY. This overrides $BUNYAN_NO_PAGER.")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "color", true, "Colorize output. Defaults to try if output stream is a TTY.")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "no-color", false, "Force no coloring.")
	RootCmd.PersistentFlags().StringVarP(&CmdOptions.Output, "output", "o", "long", "output mode/format. One of long, json, json-N, bunyan, inspect, short, simple, html, serve, server")

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
	RootCmd.PersistentFlags().StringVar(&CmdOptions.LogDestination, "log", "", "where logs are writen if given (by default, no log is generated)")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.Debug, "debug", false, "forces logging at DEBUG level")
	RootCmd.PersistentFlags().BoolVarP(&CmdOptions.Verbose, "verbose", "v", false, "runs verbosely if set")

	RootCmd.SilenceUsage = true
	cobra.OnInitialize(initConfig)
}

// initConfig reads config files and environment variable
func initConfig() {
	log := logger.Must(logger.FromContext(RootCmd.Context()))

	if len(CmdOptions.LogDestination) > 0 {
		log.ResetDestinations(CmdOptions.LogDestination)
	}
	if CmdOptions.Debug {
		log.SetFilterLevel(logger.DEBUG)
	}

	log.Infof(strings.Repeat("-", 80))
	log.Infof("Starting %s v%s (%s)", RootCmd.Name(), RootCmd.Version, runtime.GOARCH)
	log.Infof("Log Destination: %s", log)

	viper.SetConfigType("yaml")
	if len(CmdOptions.ConfigFile) > 0 { // Use config file from the flag.
		viper.SetConfigFile(CmdOptions.ConfigFile)
	} else if configDir, _ := os.UserConfigDir(); len(configDir) > 0 {
		viper.AddConfigPath(filepath.Join(configDir, "bunyan"))
		viper.SetConfigName("config.yaml")
	} else { // Old fashion way
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".bunyan")
	}

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
func runRootCommand(cmd *cobra.Command, args []string) error {
	// Here we should read from stdin or from the files
	// and pretty print the logs
	log := logger.Must(logger.FromContext(cmd.Context()))
	var scanner *bufio.Scanner

	if len(args) == 0 {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatalf("Failed to open file %s: %s", args[0], err)
			return err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	for scanner.Scan() {
		log.Infof(scanner.Text())
		var common struct {
			Time     time.Time `json:"time"`
			Level    int       `json:"level"`
			Hostname string    `json:"hostname"`
			Name     string    `json:"name"`
			PID      int       `json:"pid"`
			Message  string    `json:"msg"`
		}
		var entries map[string]interface{}
		output := strings.Builder{}
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &common); err != nil {
			log.Errorf("Failed to parse JSON: %s", err)
			fmt.Println(string(line))
			continue
		}
		if err := json.Unmarshal(line, &entries); err != nil {
			log.Errorf("Failed to parse JSON: %s", err)
			fmt.Println(string(line))
			continue
		}

		output.WriteString("[")
		output.WriteString(common.Time.Format("2006-01-02T15:04:05.000Z07:00"))
		output.WriteString("] ")
		output.WriteString(LevelColors[common.Level]) // Be sure to supprot levels not in the map
		output.WriteString(logger.Level(common.Level).String())
		output.WriteString(Reset)
		output.WriteString(": ")
		output.WriteString(common.Name)
		output.WriteString("/")
		output.WriteString(fmt.Sprintf("%d", common.PID))
		output.WriteString(" on ")
		output.WriteString(common.Hostname)
		output.WriteString(": ")
		output.WriteString(Cyan)
		output.WriteString(common.Message)
		output.WriteString(Reset)

		delete(entries, "time")
		delete(entries, "level")
		delete(entries, "hostname")
		delete(entries, "pid")
		delete(entries, "name")
		delete(entries, "msg")
		delete(entries, "v")

		output.WriteString("(")
		for key, field := range entries {
			if value, ok := field.(string); ok {
				output.WriteString(key)
				output.WriteString("=")
				output.WriteString(value)
			} else if value, ok := field.(int); ok {
				output.WriteString(key)
				output.WriteString("=")
				output.WriteString(fmt.Sprintf("%d", value))
			} else if stringer, ok := field.(fmt.Stringer); ok {
				output.WriteString(key)
				output.WriteString("=")
				output.WriteString(stringer.String())
			} else {
				output.WriteString(key)
				output.WriteString("=")
				output.WriteString(fmt.Sprintf("%v", field))
			}
			output.WriteString(", ")
		}
		output.WriteString(")")

		// If some keys are left print them one by one with \n

		fmt.Println(output.String())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read from input", err)
		return err
	}
	return nil
}
