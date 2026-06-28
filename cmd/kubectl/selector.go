package kubectl

import (
	"github.com/gildas/go-flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Selectors []Selector

type Selector struct {
	Name    string      `json:"name"    yaml:"name"`    // The name of the selector, used as a flag in the command line
	Aliases []string    `json:"aliases" yaml:"aliases"` // The aliases of the selector, used as flags in the command line
	Label   string      `json:"label"   yaml:"label"`   // The label of the selector, used to filter resources
	Usage   string      `json:"usage"   yaml:"usage"`   // The usage description of the selector
	Charts  []string    `json:"charts"  yaml:"charts"`  // The charts that support the selector
	Value   pflag.Value `json:"-"       yaml:"-"`       // The value of the selector, used to filter resources
}

var kubectlSelectors Selectors

// Register registers the selector flags to the given command
func (selector *Selector) Register(cmd *cobra.Command) {
	selector.register(cmd, selector.Name, selector.Usage)
	for _, alias := range selector.Aliases {
		selector.register(cmd, alias, selector.Usage)
	}
}

// HasFlag checks if the selector has a flag set in the command line
func (selector Selector) HasFlag(cmd *cobra.Command) (name string, ok bool) {
	if cmd.Flags().Changed(selector.Name) {
		return selector.Name, true
	}
	for _, alias := range selector.Aliases {
		if cmd.Flags().Changed(alias) {
			return alias, true
		}
	}
	return "", false
}

// HasFlag checks if any of the selectors has a flag set in the command line
func (selectors Selectors) HasFlag(cmd *cobra.Command) bool {
	for _, selector := range selectors {
		if _, found := selector.HasFlag(cmd); found {
			return true
		}
	}
	return false
}

// register registers a flag for the selector to the given command
func (selector *Selector) register(cmd *cobra.Command, name, usage string) {
	value := flags.NewEnumFlagWithFunc(cmd, "", GetResourceLabelsFunc("deployments.apps", name))
	cmd.Flags().Var(value, name, usage)
	_ = cmd.RegisterFlagCompletionFunc(value.CompletionFunc(name))
	selector.Value = value
}
