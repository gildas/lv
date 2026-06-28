# lv

`lv` is a logviewer for [Bunyan](https://github.com/trentm/node-bunyan)-based logs (like [go-logger](https://github.com/gildas/go-logger)). It also supports [pinojs](https://getpino.io) logs.

## Installation

### Linux

You can grab the latest Debian/Ubuntu, RedHat package from the [releases page](https://github.com/gildas/lv/releases) and install it with the following commands:

If you use [Homebrew](https://brew.sh), you can install `lv` with:

```bash
brew install gildas/tap/lv
```

You can also install the application with snap:

```bash
sudo snap install bunyan-logviewer
sudo snap alias bunyan-logviewer lv
```

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/bunyan-logviewer)

### macOS

You can get `lv` from [Homebrew](https://brew.sh) with:

```bash
brew install gildas/tap/lv
```

### Windows

You can install `lv` with [chocolatey](https://chocolatey.org):

```bash
choco install bunyan-logviewer
```

You can also install `lv` with [scoop](https://scoop.sh):

```bash
scoop bucket add gildas https://github.com/gildas/scoop-bucket
scoop install lv
```

### Binaries

You can download the latest version of `lv` from the [releases page](https://github.com/gildas/lv/releases).

## Usage

You can read logs from a file, from a pipe, or from a Kubernetes cluster.

```bash
lv /path/to/logfile
```

or

```bash
tail -f /path/to/logfile | lv
```

or

```bash
lv --follow /path/to/logfile
lv -f /path/to/logfile
```

By default, `lv` will display the log in a pager with colors, if the output is a terminal (you can turn off the pager with `--no-pager`). The pager is also not used when using Kubernetes logs. Any line that cannot be unmarshaled in one of the supported format will be displayed as raw text.

It will also display the time in UTC. you can display the time in local time with the `--local` flag or use any timezone of your preference with `--time xx` where `xx` is the name of the timezone, a time difference from UTC.

```bash
lv --local /path/to/logfile
lv --time America/New_York /path/to/logfile
lv --time +02:00 /path/to/logfile
lv --time -3 /path/to/logfile
```

When the `--o short` flag is used, the time is displayed in a short format.

If you use any of the Kubernetes flags, `lv` will use the `kubectl logs` command to get the logs from the Kubernetes cluster. You can use any of the `kubectl logs` flags with `lv`. For example:

```bash
lv --namespace=my-namespace --selector=app=my-app --container=my-container --follow
```

Note that `--follow` is not enough to stream logs from Kubernetes, since it can be used to stream logs from a file. You need to use other Kubernetes flags or the `--k8s` flag to tell `lv` that you want to stream logs from Kubernetes.

If the log entries contain a `topic` and a `scope` fields, `lv` will display them in color.

You can also use `lv` to filter logs by level:

```bash
lv --level=info /path/to/logfile
```

The level follows the [go-logger](https://github.com/gildas/go-logger) format. For example:

- `--level=info` will display logs of level `info` and above
- `--level=debug` will display logs of level `debug` and above
- `--level 'INFO;DEBUG{topic};TRACE{:scope}'` will display logs of level `info`and `debug` for any entry with topic `topic` and `trace` for any entry with scope `scope`.

You can also filter logs with the `--filter` flag. The filter is similar to a [JSONPath](https://goessner.net/articles/JsonPath/) expression. For example:

```bash
lv --filter '.field == "value"' /path/to/logfile
lv --filter '.field1 == .field2' /path/to/logfile
lv --filter '.field =~ /regexp/' /path/to/logfile
lv --filter '.field1 == true && .field2 == 12' /path/to/logfile
```

### Flags

Here is a list of the flags you can use with `lv`:

```txt
  --all-containers                     Get all containers' logs in the pod(s).
  --all-pods                           Get logs from all pod(s). Sets prefix to true.
  --app string                         The name of the application to use for logs
  --application string                 The name of the application to use for logs
  --as string                          Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
  --as-group stringArray               Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
  --as-uid string                      UID to impersonate for the operation.
  --as-user-extra stringArray          Key=value pairs that describe user extra fields to be impersonated for the operation. This flag can be repeated to specify multiple extra fields.
  --cache-dir string                   Default cache directory
  --certificate-authority string       Path to a cert file for the certificate authority
  --client-certificate string          Path to a client certificate file for TLS
  --client-key string                  Path to a client key file for TLS
  --cluster string                     The name of the kubeconfig cluster to use
  --color                              Colorize output always, even if the output stream is not a TTY. (default true)
  --completion string                  Generates completion script for bash, zsh, fish, or powershell
  --condition string                   Run each log message through the filter.
  --config string                      config file (default is /home/gildas/.config/logviewer/config.yaml)
  --connector string                   The name of the connector to use for logs
  -c, --container string               Print the logs of this container
  --context string                     The name of the kubeconfig context to use
  --debug                              forces logging at DEBUG level
  --disable-compression                If true, opt-out of response compression for all requests to the server
  --filter string                      Run each log message through the filter.
  -f, --follow                         Specify if the logs should be streamed
  -h, --help                           help for lv
  --ignore-errors                      If watching / following pod logs, allow for any errors that occur to be non-fatal
  --insecure-skip-tls-verify           If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
  --insecure-skip-tls-verify-backend   Skip verifying the identity of the kubelet that logs are requested from.  In theory, an attacker could provide invalid log content back. You might want to use this if your kubelet serving certificates have expired.
  -k, --key string                     Use the given key to decrypt obfuscated log entries. The key must be 16, 24, or 32 bytes long.
  --kubeconfig string                  Path to the kubeconfig file to use for CLI requests.
  --kuberc string                      Path to the kuberc file to use for preferences. This can be disabled by exporting KUBECTL_KUBERC=false feature gate or turning off the feature KUBERC=off.
  --level string                       Only shows log entries with a level at or above the given value.
  --limit-bytes int                    Maximum bytes of logs to return. Defaults to no limit.
  -L, --local                          Display time field in local time, rather than UTC.
  --log string                         where logs are writen if given (by default, no log is generated)
  --log-flush-frequency duration       Maximum number of seconds between log flushes (default 5s)
  --match-server-version               Require server version to match client version
  --max-log-requests int               Maximum number of concurrent logs to follow when using by a selector. Defaults to 5. (default 5)
  -n, --namespace string               If present, the namespace scope for this CLI request
  --no-color                           Do not colorize output. By default, the output is colorized if stdout is a TTY
  --no-pager less                      Do not pipe output into a pager. By default, the output is piped throug less (or $PAGER if set), if stdout is a TTY (default true)
  -o, --output string                  output mode/format. One of long, json, json-N, logviewer, inspect, short, simple, html, serve, server (default "long")
  --password string                    Password for basic authentication to the API server.
  --platform string                    The name of the platform to use for logs
  --pod-running-timeout duration       The length of time (like 5s, 2m, or 3h, higher than zero) to wait until at least one pod is running
  --prefix string                      Prefix each log line with the log source (pod name and container name)
  -p, --previous                       If true, print the logs for the previous instance of the container in a pod if it exists.
  --profile string                     Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex|trace)
  --profile-output string              Name of the file to write the profile to
  --provider string                    The name of the provider to use for logs
  --release string                     The name of the Helm release to use for logs
  --request-timeout duration           The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.
  --role string                        The name of the role to use for logs
  -l, --selector string                Selector (label query) to filter on, supports '=', '==', '!=', 'in', 'notin'.(e.g. -l key1=value1,key2=value2,key3 in (value3)). Matching objects must satisfy all of the specified label constraints.
  -s, --server string                  The address and port of the Kubernetes API server
  --since duration                     Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time / since may be used.
  --since-time time                    Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time / since may be used.
  --tail int                           Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided. (default -1)
  --tier string                        The name of the tier to use for logs
  --time string                        Display time field in the given timezone.
  --timestamps                         Include timestamps on each line in the log output
  --tls-server-name string             Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
  --token string                       Bearer token for authentication to the API server
  --user string                        The name of the kubeconfig user to use
  --username string                    Username for basic authentication to the API server
  -v, --verbose                        runs verbosely if set
  --version                            version for lv
  --vmodule string                     comma-separated list of pattern=N settings for file-filtered logging (only works for the default text log format)
  --warnings-as-errors                 Treat warnings received from the server as errors and exit with a non-zero exit code
```

The key must be 16, 24, or 32 bytes long.

### Environment Variables

`lv` uses the following environment variables:

- `LV_COLOR` to force colorization of the output
- `LV_LOCAL` to display the time in local time
- `LV_TIMEZONE` to display the time in a specific timezone
- `LV_OBFUSCATIONKEY` to specify the key used to decrypt obfuscated log entries

The command line flags have precedence over the environment variables.

### Configuration file

You can also configure `lv` with a configuration file. The configuration file is a YAML file called `config.yaml` and should be stored in the subfolder `logviewer` of [os.UserConfigDir](https://pkg.go.dev/os#UserConfigDir).

Here is an example of a configuration file:

```yaml
color: true
timezone: Europe/Paris
output: short
obfuscationKey: 1231213
```

Here are all the configuration options you can use in the configuration file:

- `color`: (boolean) to force colorization of the output,  
  environment variable `LV_COLOR`
- `follow`: (boolean) to follow the logs in real-time,  
  environment variable `LV_FOLLOW`
- `obfuscationKey`: (string) to specify the key used to decrypt obfuscated log entries,  
  environment variable `LV_OBFUSCATIONKEY`
- `output`: (string) to specify the output format. One of `long`, `logviewer`, `short`, `simple`, `html`, `serve`, `server`,  
  environment variable `LV_OUTPUT`
- `timezone`: (string) to display the time in a specific timezone,  
  environment variable `LV_TIMEZONE`

The environment variables and the command line flags have precedence over the configuration file.

#### Custom Selectors

You can also configure custom selectors in the configuration file. Here is an example of a configuration file with custom selectors:

```yaml
selectors:
  - name: application                # The name of the selector to use in the command line
    aliases: [app]                   # Aliases for the selector that can be used in the command line
    label: "app.kubernetes.io/name"  # The label to use for the selector in kubectl logs command.
                                     # If not specified, the name of the selector will be used.
    usage: "select application logs" # The usage of the selector to display in the command line help
```

This would be used in the command line as follows:

```bash
lv --follow --tail -1 --application=my-app
```

### Completion

`lv` supports shell completion for `bash`, `fish`, `PowerShell`, and `zsh`.

#### Bash

To enable completion, run the following command:

```bash
source <(lv --completion bash)
```

You can also add this line to your `~/.bashrc` file to enable completion for every new shell.

```bash
lv --completion bash > ~/.bashrc
```

#### Fish

To enable completion, run the following command:

```bash
lv --completion fish | source
```

You can also add this line to your `~/.config/fish/config.fish` file to enable completion for every new shell.

```bash
lv --completion fish > ~/.config/fish/completions/lv.fish
```

#### Powershell

To enable completion, run the following command:

```pwsh
lv --completion powershell | Out-String | Invoke-Expression
```

You can also add the output of the above command to your `$PROFILE` file to enable completion for every new shell.

#### zsh

To enable completion, run the following command:

```bash
source <(lv --completion zsh)
```

You can also add this line to your functions folder to enable completion for every new shell.

```bash
lv --completion zsh > "~/${fpath[1]}/_lv"
```

On macOS, you can add the completion to the brew functions:

```bash
lv --completion zsh > "$(brew --prefix)/share/zsh/site-functions/_lv"
```

## Caveats

Not all the output formats are implemented yet.

## Troubleshooting

`lv` uses [go-logger](https://github.com/gildas/go-logger) to write its own logs. You can enable the logs with the `--log` flag. By default `lv` does not log anything. `lv` will read the [go-logger](https://github.com/gildas/go-logger) environment variables with the prefix `LV_` (For example `LV_LOG_LEVEL`).

## TODO

- Add support to read logs from `aws`, `gcp` and `azure` services.
