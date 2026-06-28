package kubectl

import (
	"bytes"
	"context"
	"slices"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/common"
	"github.com/spf13/cobra"
)

var resourceTypes = map[string]flags.AllowedFunc{
	"daemonsets":             GetResourceNamesFunc("daemonsets.apps"),
	"deployments":            GetResourceNamesFunc("deployments.apps"),
	"jobs":                   GetResourceNamesFunc("jobs.batch"),
	"pods":                   GetResourceNamesFunc("pods"),
	"replicasets":            GetResourceNamesFunc("replicasets.apps"),
	"replicationcontrollers": GetResourceNamesFunc("replicationcontrollers"),
	"services":               GetResourceNamesFunc("services"),
	"statefulsets":           GetResourceNamesFunc("statefulsets.apps"),
}

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

	log.Debugf("Getting pods for completion in namespace %s with context %s, args: %s, toComplete: %s", kubectlNamespace, kubectlContext, args, toComplete)
	if len(toComplete) > 0 {
		log.Debugf("Filtering resource types for toComplete: %s", toComplete)
		resources := []string{}
		components := strings.Split(toComplete, "/")
		resourceTypeToComplete := components[0]
		resourceToComplete := ""
		if len(components) > 1 {
			resourceToComplete = components[1]
		}
		for resourceType, getter := range resourceTypes {
			if resourceType == resourceTypeToComplete {
				log.Debugf("Adding resource type %s to completion", resourceType)
				collected, err := getter(ctx, cmd, args, resourceToComplete)
				if err != nil {
					log.Errorf("Error getting resources for type %s: %v", resourceType, err)
					continue
				}
				log.Debugf("Collected resources for type %s: %s", resourceType, collected)
				resources = append(resources, core.Map(collected, func(resource string) string { return resourceType + "/" + resource })...)
				return common.FilterValidArgs(resources, args, toComplete), nil
			}
		}
		if len(resources) > 0 {
			return common.FilterValidArgs(resources, args, toComplete), nil
		}
	}
	params := []string{"get", "pods", "--context", kubectlContext, "--namespace", kubectlNamespace}
	if selector := BuildSelectorArgs(cmd); len(selector) > 0 {
		params = append(params, "--selector", selector)
	}
	params = append(params, "-o", "jsonpath={.items[*].metadata.name}")
	err = NewKubectl().Exec(ctx, params, &stdout, &stderr)
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

	if len(toComplete) == 0 {
		log.Debugf("No toComplete provided, adding resource types for completion")
		for resourceType := range resourceTypes {
			pods = append(pods, resourceType+"/")
		}
	} else {
		log.Debugf("Adding resource types for completion that match toComplete: %s", toComplete)
		for resourceType := range resourceTypes {
			if strings.HasPrefix(resourceType, toComplete) {
				log.Debugf("Adding resource type %s to completion", resourceType)
				pods = append(pods, resourceType+"/")
			}
		}
	}

	return common.FilterValidArgs(pods, args, toComplete), nil
}
