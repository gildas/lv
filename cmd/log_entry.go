package cmd

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/gildas/go-errors"
)

type LogEntry struct {
	Time     time.Time `json:"time"`
	Level    LogLevel  `json:"level"`
	Hostname string    `json:"hostname"`
	Name     string    `json:"name"`
	PID      int64     `json:"pid"`
	TaskID   int64     `json:"tid"`
	Topic    string    `json:"topic"`
	Scope    string    `json:"scope"`
	Message  string    `json:"msg"`
	Fields   map[string]any
}

func (entry LogEntry) Write(context context.Context, output io.Writer, options *OutputOptions) {
	_, _ = output.Write([]byte("["))
	entry.writeString(output, options, "[")
	if options.LocalTime {
		entry.writeString(output, options, entry.Time.Local().Format("2006-01-02T15:04:05.000"))
	} else {
		entry.writeString(output, options, entry.Time.UTC().Format("2006-01-02T15:04:05.000"))
	}
	entry.writeString(output, options, "] ")
	entry.Level.Write(output, options)
	entry.writeString(output, options, ": ")
	entry.writeString(output, options, entry.Name)
	entry.writeString(output, options, "/")
	entry.writeInt64(output, options, entry.TaskID)
	entry.writeString(output, options, " on ")
	entry.writeString(output, options, entry.Hostname)
	entry.writeString(output, options, ": ")
	if len(entry.Topic) > 0 {
		entry.writeStringWithColor(output, options, entry.Topic, Green)
		if len(entry.Scope) > 0 {
			entry.writeString(output, options, "/")
			entry.writeStringWithColor(output, options, entry.Scope, Yellow)
		}
		entry.writeString(output, options, " ")
	}
	entry.writeStringWithColor(output, options, entry.Message, Cyan)
}

func (entry LogEntry) writeString(output io.Writer, options *OutputOptions, value string) {
	_, _ = output.Write([]byte(value))
}

func (entry LogEntry) writeInt64(output io.Writer, options *OutputOptions, value int64) {
	_, _ = output.Write([]byte(strconv.FormatInt(value, 10)))
}

func (entry LogEntry) writeStringWithColor(output io.Writer, options *OutputOptions, value string, color string) {
	if options.UseColors {
		_, _ = output.Write([]byte(color))
	}
	entry.writeString(output, options, value)
	if options.UseColors {
		_, _ = output.Write([]byte(Reset))
	}
}

// UnmarshalJSON unmarshal data into this
func (entry *LogEntry) UnmarshalJSON(payload []byte) (err error) {
	var data map[string]any
	var ok bool
	var merr errors.MultiError

	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	entry.Fields = map[string]any{}
	for key, value := range data {
		switch key {
		case "hostname":
			if entry.Hostname, ok = value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("hostname", value))
			}
		case "name":
			if entry.Hostname, ok = value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("name", value))
			}
		case "topic":
			if entry.Topic, ok = value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("topic", value))
			}
		case "scope":
			if entry.Scope, ok = value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("scope", value))
			}
		case "msg":
			if entry.Message, ok = value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("msg", value))
			}
		case "level":
			if svalue, ok := value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("level", value))
			} else if ivalue, err := strconv.ParseInt(svalue, 10, 64); err != nil {
				merr.Append(errors.Join(errors.ArgumentInvalid.With("level", value), err))
			} else {
				entry.Level = LogLevel(ivalue)
			}
		case "pid":
			if svalue, ok := value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("level", value))
			} else if ivalue, err := strconv.ParseInt(svalue, 10, 64); err != nil {
				merr.Append(errors.Join(errors.ArgumentInvalid.With("pid", value), err))
			} else {
				entry.PID = ivalue
			}
		case "tid":
			if svalue, ok := value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("level", value))
			} else if ivalue, err := strconv.ParseInt(svalue, 10, 64); err != nil {
				merr.Append(errors.Join(errors.ArgumentInvalid.With("level", value), err))
			} else {
				entry.TaskID = ivalue
			}
		case "time":
			if tvalue, ok := value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("time", value))
			} else if entry.Time, err = time.Parse(time.RFC3339, tvalue); err != nil {
				merr.Append(errors.Join(errors.ArgumentInvalid.With("time", value), err))
			}
		case "v":
			// ignore
		default:
			entry.Fields[key] = value
		}
	}
	return
}
