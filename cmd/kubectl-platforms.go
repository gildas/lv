package cmd

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// KubeCtlGetPlatforms gets the platforms for the current context
func KubeCtlGetPlatforms(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "platforms")
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

	log.Debugf("Getting platforms for completion in namespace %s with context %s and args: %s", kubectlNamespace, kubectlContext, args)
	err = NewKubectl().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.labels.platform}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting platforms: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	platforms := []string{}
	for platform := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(platforms, platform) {
			platforms = append(platforms, platform)
		}
	}

	return FilterValidArgs(platforms, args, toComplete), nil
}
