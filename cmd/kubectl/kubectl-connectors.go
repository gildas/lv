package kubectl

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/common"
	"github.com/spf13/cobra"
)

// GetConnectors gets the pods for the current context
func GetConnectors(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "connectors")
	var stdout, stderr bytes.Buffer

	kubectlContext, err := GetCurrentContext(ctx, cmd)
	if err != nil {
		log.Errorf("Error getting current context: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	kubectlNamespace, err := GetCurrentNamespace(ctx, cmd, kubectlContext)
	if err != nil {
		log.Errorf("Error getting current namespace: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	log.Debugf("Getting connectors for completion in namespace %s with context %s and args: %s", kubectlNamespace, kubectlContext, args)
	err = New().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.labels.connector}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting connectors: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	connectors := []string{}
	for connector := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(connectors, connector) {
			connectors = append(connectors, connector)
		}
	}

	return common.FilterValidArgs(connectors, args, toComplete), nil
}
