package kubectl

import (
	"bytes"
	"context"
	"encoding/json"
	"slices"

	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/common"
	"github.com/spf13/cobra"
)

type HelmRelease struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Revision   string `json:"revision"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
}

// GetReleases gets the helm releases for the current context
func GetReleases(ctx context.Context, cmd *cobra.Command, args []string, toComplete string) ([]string, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("helm", "releases")
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
	err = NewHelm().Exec(ctx, []string{"list", "--kube-context", kubectlContext, "--namespace", kubectlNamespace, "-o", "json"}, &stdout, &stderr)
	if err != nil {
		log.Errorf("Error getting names: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	var helmReleases []HelmRelease
	if err := json.Unmarshal(stdout.Bytes(), &helmReleases); err != nil {
		log.Errorf("Error unmarshalling helm releases: ", err)
		log.Errorf("Stderr: %s", stderr.String())
		return nil, err
	}

	releases := []string{}
	for _, release := range helmReleases {
		if !slices.Contains(releases, release.Name) {
			releases = append(releases, release.Name)
		}
	}

	return common.FilterValidArgs(releases, args, toComplete), nil
}
