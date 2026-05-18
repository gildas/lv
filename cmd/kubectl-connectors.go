package cmd

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// KubeCtlGetConnectors gets the pods for the current context
func KubeCtlGetConnectors(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "connectors")
	var stdout, stderr bytes.Buffer

	kubectlContext, err := KubeCtlGetCurrentContext(ctx, cmd)
	if err != nil {
		log.Errorf("Error getting current context: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	kubectlNamespace, err := KubeCtlGetCurrentNamespace(ctx, cmd)
	if err != nil {
		log.Errorf("Error getting current namespace: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	log.Debugf("Getting pods for completion with args: %s", args)
	err = NewKubectl().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.labels.connector}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting pods: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	connectors := []string{}
	for connector := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(connectors, connector) {
			connectors = append(connectors, connector)
		}
	}

	return FilterValidArgs(connectors, args, toComplete), nil
}
