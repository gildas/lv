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

// GetPods gets the pods for the current context
func GetPods(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "pods")
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

	log.Debugf("Getting pods for completion in namespace %s with context %s and args: %s", kubectlNamespace, kubectlContext, args)
	err = New().Exec(ctx, []string{"get", "pods", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.name}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting pods: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	pods := []string{}
	for pod := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(pods, pod) {
			pods = append(pods, pod)
		}
	}

	return common.FilterValidArgs(pods, args, toComplete), nil
}
