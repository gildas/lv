package kubectl

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LogsOptions struct {
	AllContainers                bool
	AllPods                      bool
	As                           string
	AsGroup                      []string
	AsUID                        string
	AsUserExtra                  []string
	CacheDir                     string
	CertificateAuthority         string
	ClientCertificate            string
	ClientKey                    string
	Cluster                      string
	Container                    string
	Context                      *flags.EnumFlag
	DisableCompression           bool
	IgnoreErrors                 bool
	InsecureSkipTLSVerify        bool
	InsecureSkipTLSVerifyBackend bool
	Kubeconfig                   string
	KubeRC                       string
	LimitBytes                   int64
	LogFlushFrequency            time.Duration
	MatchServerVersion           bool
	MaxLogRequests               int
	Namespace                    *flags.EnumFlag
	Password                     string
	PodRunningTimeout            time.Duration
	Prefix                       string
	Previous                     bool
	Profile                      string
	ProfileOutput                string
	Release                      *flags.EnumFlag
	RequestTimeout               time.Duration
	Selector                     string
	Server                       string
	Since                        time.Duration
	SinceTime                    time.Time
	Tail                         int64
	Timestamps                   bool
	TLSServerName                string
	Token                        string
	User                         string
	Username                     string
	VModule                      string
	WarningsAsErrors             bool
}

var kubectlLogFlags = []string{
	"all-containers",
	"all-pods",
	"as",
	"as-group",
	"as-uid",
	"as-user-extra",
	"cache-dir",
	"certificate-authority",
	"client-certificate",
	"client-key",
	"cluster",
	"container",
	"context",
	"disable-compression",
	"ignore-errors",
	"insecure-skip-tls-verify",
	"insecure-skip-tls-verify-backend",
	"kubeconfig",
	"kuberc",
	"limit-bytes",
	"log-flush-frequency",
	"match-server-version",
	"max-log-requests",
	"namespace",
	"password",
	"pod-running-timeout",
	"prefix",
	"previous",
	"profile",
	"profile-output",
	"request-timeout",
	"selector",
	"server",
	"since",
	"since-time",
	"tail",
	"timestamps",
	"tls-server-name",
	"token",
	"user",
	"username",
	"vmodule",
	"warnings-as-errors",
}

// CreateLogsFlags creates the flags for the kubectl logs command
func CreateLogsFlags(cmd *cobra.Command) (options LogsOptions) {
	options.Context = flags.NewEnumFlagWithFunc(cmd, "", GetContexts)
	options.Namespace = flags.NewEnumFlagWithFunc(cmd, "", GetNamespaces)
	options.Release = flags.NewEnumFlagWithFunc(cmd, "", GetReleases)

	cmd.Flags().BoolVar(&options.AllContainers, "all-containers", false, "Get all containers' logs in the pod(s).")
	cmd.Flags().BoolVar(&options.AllPods, "all-pods", false, "Get logs from all pod(s). Sets prefix to true.")
	cmd.Flags().StringVar(&options.As, "as", "", "Username to impersonate for the operation. User could be a regular user or a service account in a namespace.")
	cmd.Flags().StringArrayVar(&options.AsGroup, "as-group", []string{}, "Group to impersonate for the operation, this flag can be repeated to specify multiple groups.")
	cmd.Flags().StringVar(&options.AsUID, "as-uid", "", "UID to impersonate for the operation.")
	cmd.Flags().StringArrayVar(&options.AsUserExtra, "as-user-extra", []string{}, "Key=value pairs that describe user extra fields to be impersonated for the operation. This flag can be repeated to specify multiple extra fields.")
	cmd.Flags().StringVar(&options.CacheDir, "cache-dir", "", "Default cache directory")
	cmd.Flags().StringVar(&options.CertificateAuthority, "certificate-authority", "", "Path to a cert file for the certificate authority")
	cmd.Flags().StringVar(&options.ClientCertificate, "client-certificate", "", "Path to a client certificate file for TLS")
	cmd.Flags().StringVar(&options.ClientKey, "client-key", "", "Path to a client key file for TLS")
	cmd.Flags().StringVar(&options.Cluster, "cluster", "", "The name of the kubeconfig cluster to use")
	cmd.Flags().StringVarP(&options.Container, "container", "c", "", "Print the logs of this container")
	cmd.Flags().Var(options.Context, "context", "The name of the kubeconfig context to use")
	cmd.Flags().BoolVar(&options.DisableCompression, "disable-compression", false, "If true, opt-out of response compression for all requests to the server")
	cmd.Flags().BoolVar(&options.IgnoreErrors, "ignore-errors", false, "If watching / following pod logs, allow for any errors that occur to be non-fatal")
	cmd.Flags().BoolVar(&options.InsecureSkipTLSVerify, "insecure-skip-tls-verify", false, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	cmd.Flags().BoolVar(&options.InsecureSkipTLSVerifyBackend, "insecure-skip-tls-verify-backend", false, "Skip verifying the identity of the kubelet that logs are requested from.  In theory, an attacker could provide invalid log content back. You might want to use this if your kubelet serving certificates have expired.")
	cmd.Flags().StringVar(&options.Kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests.")
	cmd.Flags().StringVar(&options.KubeRC, "kuberc", "", "Path to the kuberc file to use for preferences. This can be disabled by exporting KUBECTL_KUBERC=false feature gate or turning off the feature KUBERC=off.")
	cmd.Flags().Int64Var(&options.LimitBytes, "limit-bytes", 0, "Maximum bytes of logs to return. Defaults to no limit.")
	cmd.Flags().DurationVar(&options.LogFlushFrequency, "log-flush-frequency", 5*time.Second, "Maximum number of seconds between log flushes")
	cmd.Flags().BoolVar(&options.MatchServerVersion, "match-server-version", false, "Require server version to match client version")
	cmd.Flags().IntVar(&options.MaxLogRequests, "max-log-requests", 5, "Maximum number of concurrent logs to follow when using by a selector. Defaults to 5.")
	cmd.Flags().VarP(options.Namespace, "namespace", "n", "If present, the namespace scope for this CLI request")
	cmd.Flags().StringVar(&options.Password, "password", "", "Password for basic authentication to the API server.")
	cmd.Flags().DurationVar(&options.PodRunningTimeout, "pod-running-timeout", 0, "The length of time (like 5s, 2m, or 3h, higher than zero) to wait until at least one pod is running")
	cmd.Flags().StringVar(&options.Prefix, "prefix", "", "Prefix each log line with the log source (pod name and container name)")
	cmd.Flags().BoolVarP(&options.Previous, "previous", "p", false, "If true, print the logs for the previous instance of the container in a pod if it exists.")
	cmd.Flags().StringVar(&options.Profile, "profile", "", "Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex|trace)")
	cmd.Flags().StringVar(&options.ProfileOutput, "profile-output", "", "Name of the file to write the profile to")
	cmd.Flags().DurationVar(&options.RequestTimeout, "request-timeout", 0, "The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.")
	cmd.Flags().StringVarP(&options.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', '!=', 'in', 'notin'.(e.g. -l key1=value1,key2=value2,key3 in (value3)). Matching objects must satisfy all of the specified label constraints.")
	cmd.Flags().StringVarP(&options.Server, "server", "s", "", "The address and port of the Kubernetes API server")
	cmd.Flags().DurationVar(&options.Since, "since", 0, "Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time / since may be used.")
	cmd.Flags().TimeVar(&options.SinceTime, "since-time", time.Time{}, []string{time.RFC3339}, "Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time / since may be used.")
	cmd.Flags().Int64Var(&options.Tail, "tail", -1, "Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided.")
	cmd.Flags().BoolVar(&options.Timestamps, "timestamps", false, "Include timestamps on each line in the log output")
	cmd.Flags().StringVar(&options.TLSServerName, "tls-server-name", "", "Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used")
	cmd.Flags().StringVar(&options.Token, "token", "", "Bearer token for authentication to the API server")
	cmd.Flags().StringVar(&options.User, "user", "", "The name of the kubeconfig user to use")
	cmd.Flags().StringVar(&options.Username, "username", "", "Username for basic authentication to the API server")
	cmd.Flags().StringVar(&options.VModule, "vmodule", "", "comma-separated list of pattern=N settings for file-filtered logging (only works for the default text log format)")
	cmd.Flags().BoolVar(&options.WarningsAsErrors, "warnings-as-errors", false, "Treat warnings received from the server as errors and exit with a non-zero exit code")
	if IsHelmAvailable() {
		cmd.Flags().Var(options.Release, "release", "The name of the Helm release to use for logs")
	}

	_ = cmd.RegisterFlagCompletionFunc(options.Context.CompletionFunc("context"))
	_ = cmd.RegisterFlagCompletionFunc(options.Namespace.CompletionFunc("namespace"))
	_ = cmd.RegisterFlagCompletionFunc(options.Release.CompletionFunc("release"))
	return
}

// HasLogsFlags checks if any of the kubectl logs flags or extra logs flags are present in the command
func HasLogsFlags(cmd *cobra.Command) bool {
	if cmd.Flag("k8s").Changed {
		return true
	}
	if slices.ContainsFunc(kubectlLogFlags, func(flag string) bool { return cmd.Flags().Changed(flag) }) {
		return true
	}
	return kubectlSelectors.HasFlag(cmd)
}

// BuildLogsParameters builds the parameters for the kubectl logs command based on the flags present in the command
func BuildLogsParameters(cmd *cobra.Command) (params []string) {
	params = []string{"logs"}
	// --follow is common for both kubectl and files, so we need to add it here
	if viper.GetBool("follow") {
		params = append(params, "--follow")
	}
	for _, flag := range kubectlLogFlags {
		if cmd.Flags().Changed(flag) {
			params = append(params, "--"+flag)
			// If the flag has a value, we need to add it as well
			if cmd.Flags().Lookup(flag).Value.String() != "" && cmd.Flags().Lookup(flag).Value.Type() != "bool" {
				params = append(params, cmd.Flags().Lookup(flag).Value.String())
			}
		}
	}
	if selector := BuildSelectorArgs(cmd); len(selector) > 0 {
		params = append(params, "--selector", selector)
	}
	return
}

// BuildSelectorArgs builds the Kubernetes selector
func BuildSelectorArgs(cmd *cobra.Command) string {
	log := logger.Must(logger.FromContext(cmd.Context()))

	selectors := []string{}
	for _, selector := range kubectlSelectors {
		if name, found := selector.HasFlag(cmd); found {
			log.Debugf("Selector %s found with flag %s", selector.Name, name)
			// If the flag has a value, we need to add it as well
			if cmd.Flags().Lookup(name).Value.String() != "" && cmd.Flags().Lookup(name).Value.Type() != "bool" {
				selectors = append(selectors, fmt.Sprintf("%s=%s", selector.GetLabel(), cmd.Flags().Lookup(name).Value.String()))
			} else {
				selectors = append(selectors, selector.GetLabel())
			}
		}
	}
	return strings.Join(selectors, ",")
}
