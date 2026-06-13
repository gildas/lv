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

// GetNames gets the pods for the current context
func GetNames(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "names")
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

	log.Debugf("Getting names for completion with args: %s", args)
	err = NewKubectl().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.name}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting names: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	names := []string{}
	for name := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(names, name) {
			names = append(names, name)
		}
	}

	return common.FilterValidArgs(names, args, toComplete), nil
}
