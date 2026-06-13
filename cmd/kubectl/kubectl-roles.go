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

// GetRoles gets the roles for the current context
func GetRoles(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "roles")
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

	log.Debugf("Getting roles for completion in namespace %s with context %s and args: %s", kubectlNamespace, kubectlContext, args)
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

	return common.FilterValidArgs(roles, args, toComplete), nil
}
