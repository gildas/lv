package cmd

import (
	"bytes"
	"context"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// KubeCtlGetContexts gets the contexts for the current kubeconfig
func KubeCtlGetContexts(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "contexts")
	var stdout, stderr bytes.Buffer

	log.Debugf("Getting contexts for completion with args: %s", args)
	err := NewKubectl().Exec(ctx, []string{"config", "get-contexts", "-o", "name"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting contexts: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	contexts := []string{}
	for _, context := range strings.Fields(stdout.String()) {
		contexts = append(contexts, context)
	}

	return FilterValidArgs(contexts, args, toComplete), nil
}

// KubeCtlGetCurrentContext gets the current context for the current kubeconfig
func KubeCtlGetCurrentContext(ctx context.Context, cmd *cobra.Command) (string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "current-context")

	if cmd.Flags().Changed("context") {
		return CmdOptions.Context.Value, nil
	}

	var stdout, stderr bytes.Buffer

	log.Debugf("Getting current context")
	err := NewKubectl().Exec(ctx, []string{"config", "current-context"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting current context: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
