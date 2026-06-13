package cmd

import (
	"bufio"
	"context"
	"crypto/aes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/kubectl"
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
	kubectl.LogsOptions
	kubectl.ExtraLogsOptions
	Completion     *flags.EnumFlag
	ConfigFile     string
	CipherKey      string
	LogDestination string
	Timezone       string
	UsePager       bool
	Verbose        bool
	Debug          bool
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Short:             "pretty-print logviewer logs from stdin, file(s), or Kubernetes resources",
	Long:              "logviewer is a simple and fast JSON log viewer. It reads log entries from given files, stdin, or Kubernetes resources and pretty-prints them to stdout.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: validRootArgs,
	RunE:              runRootCommand,
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
	RootCmd.PersistentFlags().StringVar(&CmdOptions.Filter, "filter", "", "Run each log message through the filter.")
	RootCmd.PersistentFlags().StringVar(&CmdOptions.Filter, "condition", "", "Run each log message through the filter.")
	RootCmd.PersistentFlags().StringVarP(&CmdOptions.CipherKey, "key", "k", "", "Use the given key to decrypt obfuscated log entries. The key must be 16, 24, or 32 bytes long.")
	RootCmd.PersistentFlags().BoolP("local", "L", false, "Display time field in local time, rather than UTC.")
	RootCmd.PersistentFlags().StringVar(&CmdOptions.Timezone, "time", "", "Display time field in the given timezone.")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UsePager, "no-pager", true, "Do not pipe output into a pager. By default, the output is piped throug `less` (or $PAGER if set), if stdout is a TTY")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "no-color", false, "Do not colorize output. By default, the output is colorized if stdout is a TTY")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.UseColors, "color", true, "Colorize output always, even if the output stream is not a TTY.")
	RootCmd.PersistentFlags().VarP(CmdOptions.Output, "output", "o", "output mode/format. One of long, json, json-N, logviewer, inspect, short, simple, html, serve, server")
	RootCmd.PersistentFlags().StringVar(&CmdOptions.LogDestination, "log", "", "where logs are writen if given (by default, no log is generated)")
	RootCmd.PersistentFlags().BoolVar(&CmdOptions.Debug, "debug", false, "forces logging at DEBUG level")
	RootCmd.PersistentFlags().BoolVarP(&CmdOptions.Verbose, "verbose", "v", false, "runs verbosely if set")
	CmdOptions.LogsOptions = kubectl.CreateLogsFlags(RootCmd)
	CmdOptions.ExtraLogsOptions = kubectl.CreateSelectorFlags(RootCmd)

	_ = RootCmd.RegisterFlagCompletionFunc(CmdOptions.Output.CompletionFunc("output"))
	_ = RootCmd.RegisterFlagCompletionFunc(CmdOptions.Completion.CompletionFunc("completion"))

	RootCmd.SilenceUsage = true // Do not show usage when an error occurs
	cobra.OnInitialize(func() {
		if err := Initialize(RootCmd); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize: %s\n", err)
			os.Exit(1)
		}
	})
}

// validRootArgs validates the arguments passed to the root command
func validRootArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child("completion", "validate-args")

	// If the command flags indicate we are using Kubernetes resources, we should complete pods
	if kubectl.HasLogsFlags(cmd) {
		log.Debugf("Kubectl logs flags detected, completing pods")
		if pods, err := kubectl.GetPods(cmd.Context(), cmd, args, toComplete); err == nil {
			return pods, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	// If not, we should complete files in the current directory
	log.Debugf("Letting the shell to complete the files in the current directory for args: %s, toComplete: %s", args, toComplete)
	return nil, cobra.ShellCompDirectiveDefault
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
	CmdOptions.UsePager = isStdoutTTY() && isStdinTTY() && !kubectl.HasLogsFlags(cmd)
	if cmd.Flags().Changed("no-pager") || viper.GetBool("no-pager") {
		CmdOptions.UsePager = false
	}
	CmdOptions.Output.Value = viper.GetString("output")

	if len(viper.GetString("obfuscationKey")) > 0 {
		cipherBlock, err := aes.NewCipher([]byte(viper.GetString("obfuscationKey")))
		if err != nil {
			log.Fatalf("Failed to create cipher block: %s", err)
			return err
		}
		log.SetObfuscationKey(cipherBlock)
	}

	if viper.GetBool("local") {
		log.Infof("Displaying local time")
		CmdOptions.Location = time.Local
	} else if CmdOptions.Location, err = ParseLocation(viper.GetString("timezone")); err != nil {
		log.Fatalf("Failed to load timezone %s: %s", viper.GetString("timezone"), err)
		return err
	}
	log.Infof("Displaying time at location: %s", CmdOptions.Location)

	// If some of the Kubectl Logs flags are set, we should execute kubectl logs command and read from its output
	if kubectl.HasLogsFlags(cmd) {
		pipeReader, pipeWriter, err := os.Pipe()
		if err != nil {
			log.Fatalf("Failed to create pipe: %s", err)
			return err
		}
		reader = bufio.NewReader(pipeReader)

		log.Infof("Executing kubectl logs command with the given flags")
		go func() {
			defer func() { _ = pipeWriter.Close() }()
			params := kubectl.BuildLogsParameters(cmd)
			params = append(params, args...)
			if err := kubectl.NewKubectl().Exec(cmd.Context(), params, pipeWriter, pipeWriter); err != nil {
				log.Fatalf("Failed to execute kubectl logs command: %s", err)
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}()
	} else if len(args) == 0 {
		log.Infof("Reading from stdin")
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatalf("Failed to open file %s: %s", args[0], err)
			return err
		}
		defer func() { _ = file.Close() }()
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
	var filter = filters.AsFilter()

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
			_, _ = fmt.Fprintln(outstream, string(line))
			continue
		}
		if filter.Filter(cmd.Context(), entry) {
			output := strings.Builder{}

			entry.Write(cmd.Context(), &output, &CmdOptions.OutputOptions)
			if output.Len() > 0 {
				_, _ = fmt.Fprintln(outstream, output.String())
			}
		}
	}
	if err != nil && !errors.Is(err, io.EOF) {
		log.Fatalf("Failed to read from input", err)
		return err
	}
	return nil
}
