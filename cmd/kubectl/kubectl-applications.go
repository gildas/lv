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

// GetApplications gets the applications for the current context
func GetApplications(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "applications")
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

	log.Debugf("Getting applications for completion in namespace %s with context %s and args: %s", kubectlNamespace, kubectlContext, args)
	err = New().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.labels.app\\.kubernetes\\.io/name}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting applications: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	applications := []string{}
	for application := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(applications, application) {
			applications = append(applications, application)
		}
	}

	return common.FilterValidArgs(applications, args, toComplete), nil
}
