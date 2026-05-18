package cmd

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// KubeCtlGetRoles gets the roles for the current context
func KubeCtlGetRoles(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "roles")
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

	log.Debugf("Getting roles for completion with args: %s", args)
	err = NewKubectl().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.labels.role}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting roles: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	roles := []string{}
	for role := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(roles, role) {
			roles = append(roles, role)
		}
	}

	return FilterValidArgs(roles, args, toComplete), nil
}
