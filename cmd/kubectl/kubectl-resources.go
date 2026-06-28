package kubectl

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/common"
	"github.com/spf13/cobra"
)

// GetResourceNamesFunc gets the kubernetes resources for the current context
func GetResourceNamesFunc(resourceType string) flags.AllowedFunc {
	return GetResourceLabelsFunc(resourceType, "name")
}

// GetResourceLabelsFunc gets the kubernetes resources for the current context with a specific label selector
func GetResourceLabelsFunc(resourceType string, labelSelector string) flags.AllowedFunc {
	if labelSelector != "name" {
		labelSelector = "labels." + strings.ReplaceAll(labelSelector, ".", "\\.")
	}
	return func(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
		log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "resources", "resourceType", resourceType, "labelSelector", labelSelector)
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

		log.Debugf("Getting %s for completion in namespace %s with context %s, label selector %s and args: %s", resourceType, kubectlNamespace, kubectlContext, labelSelector, args)
		err = NewKubectl().Exec(ctx, []string{"get", resourceType, "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata." + labelSelector + "}"}, &stdout, &stderr)
		if err != nil {
			log.Errorf("Error getting %s: ", resourceType, err)
			log.Errorf("Stderr: %s", stderr.String())
			return nil, err
		}
		log.Debugf("Raw resources output: %s", stdout.String())

		resources := []string{}
		for resource := range strings.FieldsSeq(stdout.String()) {
			if !slices.Contains(resources, resource) {
				resources = append(resources, resource)
			}
		}

		return common.FilterValidArgs(resources, args, toComplete), nil
	}
}
