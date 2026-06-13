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

// GetResourcesFunc gets the kubernetes resources for the current context
func GetResourcesFunc(resourceType string) flags.AllowedFunc {
	return func(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
		log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "resources", "resourceType", resourceType)
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

		log.Debugf("Getting %s for completion in namespace %s with context %s and args: %s", resourceType, kubectlNamespace, kubectlContext, args)
		err = NewKubectl().Exec(ctx, []string{"get", resourceType, "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath='{.items[*].metadata.name}'"}, &stdout, &stderr)
		if err != nil {
			log.Errorf("Error getting %s: ", resourceType, err)
			log.Errorf("Stderr: %s", stderr.String())
			return nil, err
		}

		resources := []string{}
		for resource := range strings.FieldsSeq(stdout.String()) {
			if !slices.Contains(resources, resource) {
				resources = append(resources, resource)
			}
		}

		return common.FilterValidArgs(resources, args, toComplete), nil
	}
}
