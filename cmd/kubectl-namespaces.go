package cmd

import (
	"bytes"
	"context"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// KubeCtlGetNamespaces gets the namespaces for the current context
func KubeCtlGetNamespaces(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "namespaces")
	var stdout, stderr bytes.Buffer

	kubectlContext, err := KubeCtlGetCurrentContext(ctx, cmd)
	if err != nil {
		log.Errorf("Error getting current context: ", err)
		return nil, err
	}

	log.Debugf("Getting namespaces for completion with args: %s", args)
	err = NewKubectl().Exec(ctx, []string{"get", "namespaces", "--context", kubectlContext, "-o", "jsonpath={.items[*].metadata.name}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting namespaces: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	namespaces := []string{}
	for ns := range strings.FieldsSeq(stdout.String()) {
		namespaces = append(namespaces, ns)
	}

	return FilterValidArgs(namespaces, args, toComplete), nil
}

// KubeCtlGetCurrentNamespace gets the current context for the current kubeconfig
func KubeCtlGetCurrentNamespace(ctx context.Context, cmd *cobra.Command, kubectlContext string) (string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "current-namespace")

	if cmd.Flags().Changed("namespace") {
		return CmdOptions.Namespace.Value, nil
	}

	if len(kubectlContext) == 0 {
		var err error

		kubectlContext, err = KubeCtlGetCurrentContext(ctx, cmd)
		if err != nil {
			log.Errorf("Error getting current context: ", err)
			return "", err
		}
	}

	var stdout, stderr bytes.Buffer

	log.Debugf("Getting current namespace for context %s", kubectlContext)
	err := NewKubectl().Exec(ctx, []string{"config", "view", "--context", kubectlContext, "--minify", "--output", "jsonpath={..namespace}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting current namespace: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
