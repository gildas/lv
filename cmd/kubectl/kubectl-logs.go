package kubectl

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gildas/go-flags"
	"github.com/spf13/cobra"
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
	Follow                       bool
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

type ExtraLogsOptions struct {
	Connector   *flags.EnumFlag
	Platform    *flags.EnumFlag
	Provider    *flags.EnumFlag
	Release     *flags.EnumFlag
	Role        *flags.EnumFlag
	Tier        *flags.EnumFlag
	Application *flags.EnumFlag
}

var kubectlLogsFlVags = []string{
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
	"follow",
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

var kubectlExtraLogsFlVags = []string{
	"connector",
	"platform",
	"provider",
	"role",
	"tier",
	"release",
	"application",
	"app",
}

var kubectlExtraLogsSelectors = map[string]string{
	"application": "app.kubernetes.io/name",
	"app":         "app.kubernetes.io/name",
	"connector":   "connector",
	"platform":    "platform",
	"provider":    "provider",
	"role":        "role",
	"tier":        "tier",
	"release":     "app.kubernetes.io/instance",
}

// CreateLogsFlags creates the flags for the kubectl logs command
func CreateLogsFlags(cmd *cobra.Command) (options LogsOptions) {
	options.Context = flags.NewEnumFlagWithFunc(cmd, "", GetContexts)
	options.Namespace = flags.NewEnumFlagWithFunc(cmd, "", GetNamespaces)

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
	cmd.Flags().BoolVarP(&options.Follow, "follow", "f", false, "Specify if the logs should be streamed")
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

	_ = cmd.RegisterFlagCompletionFunc(options.Context.CompletionFunc("context"))
	_ = cmd.RegisterFlagCompletionFunc(options.Namespace.CompletionFunc("namespace"))
	return
}

// CreateExtraLogsFlags creates the extra flags for the kubectl logs command
//
// # These flags are helpers to build kubectl selectors to select the pods to get logs
//
// Caveat: these flags are based on the Kubernetes clusters I typically build. It would be nice to make this configurable
func CreateExtraLogsFlags(cmd *cobra.Command) (options ExtraLogsOptions) {
	options.Connector = flags.NewEnumFlagWithFunc(cmd, "", GetConnectors)
	options.Platform = flags.NewEnumFlagWithFunc(cmd, "", GetPlatforms)
	options.Provider = flags.NewEnumFlagWithFunc(cmd, "", GetProviders)
	options.Release = flags.NewEnumFlagWithFunc(cmd, "", GetReleases)
	options.Role = flags.NewEnumFlagWithFunc(cmd, "", GetRoles)
	options.Tier = flags.NewEnumFlagWithFunc(cmd, "", GetTiers)
	options.Application = flags.NewEnumFlagWithFunc(cmd, "", GetApplications)

	cmd.Flags().Var(options.Connector, "connector", "The name of the connector to use for logs")
	cmd.Flags().Var(options.Platform, "platform", "The name of the platform to use for logs")
	cmd.Flags().Var(options.Provider, "provider", "The name of the provider to use for logs")
	cmd.Flags().Var(options.Role, "role", "The name of the role to use for logs")
	cmd.Flags().Var(options.Tier, "tier", "The name of the tier to use for logs")
	cmd.Flags().Var(options.Application, "application", "The name of the application to use for logs")
	cmd.Flags().Var(options.Application, "app", "The name of the application to use for logs")
	if IsHelmAvailable() {
		cmd.Flags().Var(options.Release, "release", "The name of the Helm release to use for logs")
	}

	cmd.MarkFlagsMutuallyExclusive("connector", "platform", "provider", "application", "app", "role", "tier")

	_ = cmd.RegisterFlagCompletionFunc(options.Connector.CompletionFunc("connector"))
	_ = cmd.RegisterFlagCompletionFunc(options.Platform.CompletionFunc("platform"))
	_ = cmd.RegisterFlagCompletionFunc(options.Provider.CompletionFunc("provider"))
	_ = cmd.RegisterFlagCompletionFunc(options.Release.CompletionFunc("release"))
	_ = cmd.RegisterFlagCompletionFunc(options.Role.CompletionFunc("role"))
	_ = cmd.RegisterFlagCompletionFunc(options.Tier.CompletionFunc("tier"))
	_ = cmd.RegisterFlagCompletionFunc(options.Application.CompletionFunc("application"))
	_ = cmd.RegisterFlagCompletionFunc(options.Application.CompletionFunc("app"))
	return
}

// HasLogsFlags checks if any of the kubectl logs flags or extra logs flags are present in the command
func HasLogsFlags(cmd *cobra.Command) bool {
	return slices.ContainsFunc(kubectlLogsFlVags, func(flag string) bool {
		return cmd.Flags().Changed(flag)
	}) || slices.ContainsFunc(kubectlExtraLogsFlVags, func(flag string) bool {
		return cmd.Flags().Changed(flag)
	})
}

// BuildLogsParameters builds the parameters for the kubectl logs command based on the flags present in the command
func BuildLogsParameters(cmd *cobra.Command) (params []string) {
	params = []string{"logs"}
	for _, flag := range kubectlLogsFlVags {
		if cmd.Flags().Changed(flag) {
			params = append(params, "--"+flag)
			// If the flag has a value, we need to add it as well
			if cmd.Flags().Lookup(flag).Value.String() != "" && cmd.Flags().Lookup(flag).Value.Type() != "bool" {
				params = append(params, cmd.Flags().Lookup(flag).Value.String())
			}
		}
	}
	selectors := []string{}
	for _, flag := range kubectlExtraLogsFlVags {
		if cmd.Flags().Changed(flag) {
			// If the flag has a value, we need to add it as well
			if cmd.Flags().Lookup(flag).Value.String() != "" && cmd.Flags().Lookup(flag).Value.Type() != "bool" {
				selectors = append(selectors, fmt.Sprintf("%s=%s", kubectlExtraLogsSelectors[flag], cmd.Flags().Lookup(flag).Value.String()))
			}
		}
	}
	if len(selectors) > 0 {
		params = append(params, "--selector", strings.Join(selectors, ","))
	}
	return
}
