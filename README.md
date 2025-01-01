# lv

`lv` is a logviewer for [Bunyan](https://github.com/trentm/node-bunyan)-based logs (like [go-logger](https://github.com/gildas/go-logger)). It also supports [pinojs](https://getpino.io) logs.

## Installation

### Linux

You can grab the latest Debian/Ubuntu, RedHat package from the [releases page](https://github.com/gildas/lv/releases) and install it with the following commands:

If you use [Homebrew](https://brew.sh), you can install `lv` with:

```bash
brew install gildas/tap/lv
```

### macOS

You can get `lv` from [Homebrew](https://brew.sh) with:

```bash
brew install gildas/tap/lv
```

### Binaries

You can download the latest version of `lv` from the [releases page](https://github.com/gildas/lv/releases).

## Usage

You can read logs from a file or from a pipe.

```bash
lv /path/to/logfile
```

or

```bash
tail -f /path/to/logfile | lv
```

By default, `lv` will display the log in a pager with colors, if the output is a terminal (you can turn off the pager with `--no-pager`). Any line that does cannot be unmarshaled in one of the supported format will be displayed as raw text.

It will also display the time in UTC. you can display the time in local time with the `--local` flag or use any timezone of your preference with `--time xx` where `xx` is the name of the timezone, a time difference from UTC.

```bash
lv --local /path/to/logfile
lv --time America/New_York /path/to/logfile
lv --time +02:00 /path/to/logfile
lv --time -3 /path/to/logfile
```

When the `--o short` flag is used, the time is displayed in a short format.

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
      --color              Colorize output always, even if the output stream is not a TTY. (default true)
  -c, --condition string   Run each log message through the filter.
      --debug              forces lv's logging at DEBUG level
  -f, --filter string      Run each log message through the filter.
  -h, --help               help for lv
      --level string       Only shows log entries with a level at or above the given value.
  -L, --local              Display time field in local time, rather than UTC.
      --log string         where lv's logs are writen if given (by default, no log is generated)
      --no-color           Do not colorize output. By default, the output is colorized if stdout is a TTY
      --no-pager           Do not pipe output into a pager. By default, the output is piped throug less 
                           (or $PAGER if set), if stdout is a TTY (default true)
  -o, --output string      output mode/format. One of long, json, short, html, serve, server (default "long")
      --time string        Display time field in the given timezone.
  -v, --verbose            runs verbosely if set
      --version            version for lv
```

### Environment Variables

`lv` uses the following environment variables:

- `LOGVIEWER_COLOR` to force colorization of the output
- `LOGVIEWER_LOCAL` to display the time in local time
- `LOGVIEWER_TIMEZONE` to display the time in a specific timezone

The command line flags have precedence over the environment variables.

### Configuration file

You can also configure `lv` with a configuration file. The configuration file is a YAML file called `config.yaml` and should be stored in the subfolder `logviewer` of [os.UserConfigDir](https://pkg.go.dev/os#UserConfigDir).

Here is an example of a configuration file:

```yaml
color: true
local: true
output: short
```

The environment variables and the command line flags have precedence over the configuration file.

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

`lv` uses [go-logger](https://github.com/gildas/go-logger) to write its own logs. You can enable the logs with the `--log` flag. By default `lv` does not log anything.

## TODO

- Add support to read logs from `aws`, `gcp` and `azure` services.
- Add support for `k8s` logs.
