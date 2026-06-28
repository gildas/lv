package kubectl

import (
	"github.com/gildas/go-flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

// InitializeSelectors initializes the selectors by unmarshalling the configuration and registering them to the command
func InitializeSelectors(cmd *cobra.Command) error {
	if err := viper.UnmarshalKey("selectors", &kubectlSelectors); err != nil {
		return err
	}
	for i := range kubectlSelectors {
		kubectlSelectors[i].Register(cmd)
	}
	return nil
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

// Register registers the selector flags to the given command
func (selector *Selector) Register(cmd *cobra.Command) {
	selector.register(cmd, selector.Name)
	for _, alias := range selector.Aliases {
		selector.register(cmd, alias)
	}
}

// GetLabel returns the label of the selector, used to filter resources
func (selector Selector) GetLabel() string {
	if len(selector.Label) > 0 {
		return selector.Label
	}
	return selector.Name
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

// register registers a flag for the selector to the given command
func (selector *Selector) register(cmd *cobra.Command, name string) {
	value := flags.NewEnumFlagWithFunc(cmd, "", GetResourceLabelsFunc("deployments.apps", selector.GetLabel()))
	if cmd.Flags().Lookup(name) == nil {
		cmd.Flags().Var(value, name, selector.Usage)
	}
	_ = cmd.RegisterFlagCompletionFunc(value.CompletionFunc(name))
	selector.Value = value
}
