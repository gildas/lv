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

// GetTiers gets the tiers for the current context
func GetTiers(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "tiers")
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

	log.Debugf("Getting tiers for completion in namespace %s with context %s and args: %s", kubectlNamespace, kubectlContext, args)
	err = NewKubectl().Exec(ctx, []string{"get", "deployments.apps", "--context", kubectlContext, "--namespace", kubectlNamespace, "-o", "jsonpath={.items[*].metadata.labels.tier}"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting tiers: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	tiers := []string{}
	for tier := range strings.FieldsSeq(stdout.String()) {
		if !slices.Contains(tiers, tier) {
			tiers = append(tiers, tier)
		}
	}

	return common.FilterValidArgs(tiers, args, toComplete), nil
}
