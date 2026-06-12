package kubectl

import (
	"bytes"
	"context"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/common"
	"github.com/spf13/cobra"
)

// GetContexts gets the contexts for the current kubeconfig
func GetContexts(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "contexts")
	var stdout, stderr bytes.Buffer

	log.Debugf("Getting contexts for completion with args: %s", args)
	err := New().Exec(ctx, []string{"config", "get-contexts", "-o", "name"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting contexts: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	contexts := []string{}
	for context := range strings.FieldsSeq(stdout.String()) {
		contexts = append(contexts, context)
	}

	return common.FilterValidArgs(contexts, args, toComplete), nil
}

// GetCurrentContext gets the current context for the current kubeconfig
func GetCurrentContext(ctx context.Context, cmd *cobra.Command) (string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "current-context")

	if cmd.Flags().Changed("context") {
		return cmd.Flags().Lookup("context").Value.String(), nil
	}

	var stdout, stderr bytes.Buffer

	log.Debugf("Getting current context")
	err := New().Exec(ctx, []string{"config", "current-context"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting current context: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
