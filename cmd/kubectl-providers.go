package cmd

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// KubeCtlGetProviders gets the pods for the current context
func KubeCtlGetProviders(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "providers")
	var stdout, stderr bytes.Buffer

	kubectlContext, err := KubeCtlGetCurrentContext(ctx, cmd)
	if err != nil {
		log.Errorf("Error getting current context: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	kubectlNamespace, err := KubeCtlGetCurrentNamespace(ctx, cmd, kubectlContext)
	if err != nil {
		log.Errorf("Error getting current namespace: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	log.Debugf("Getting providers for completion in namespace %s with context %s and args: %s", kubectlNamespace, kubectlContext, args)
	err = NewKubectl().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.labels.provider}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting providers: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	providers := []string{}
	for provider := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(providers, provider) {
			providers = append(providers, provider)
		}
	}

	return FilterValidArgs(providers, args, toComplete), nil
}
