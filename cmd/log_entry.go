package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

// LogEntry represents a log entry
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
	Blobs    map[string]any
}

// GetField retrieves the value of a specific field from the LogEntry.
func (entry LogEntry) GetField(name string) string {
	if value, ok := entry.Fields[name]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	if name == "level" {
		return entry.Level.String()
	}
	if name == "hostname" {
		return entry.Hostname
	}
	if name == "name" {
		return entry.Name
	}
	if name == "pid" {
		return strconv.FormatInt(entry.PID, 10)
	}
	if name == "tid" {
		return strconv.FormatInt(entry.TaskID, 10)
	}
	if name == "topic" {
		return entry.Topic
	}
	if name == "scope" {
		return entry.Scope
	}
	if name == "msg" {
		return entry.Message
	}
	return ""
}

// Write writes the LogEntry to the given io.Writer output
// The output will be formatted according to the OutputOptions
func (entry LogEntry) Write(context context.Context, output io.Writer, options *OutputOptions) {
	log := logger.Must(logger.FromContext(context))

	entry.writeHeader(output, options)
	entry.writeString(output, options, ": ")
	entry.writeTopicAndScope(output, options)
	entry.writeStringWithColor(output, options, entry.Message, Cyan)

	log.Debugf("Fields: %v", entry.Fields)
	entry.writeString(output, options, " (")
	if len(entry.Fields) > 0 {
		index := 0
		for key, field := range entry.Fields {
			if index > 0 {
				entry.writeString(output, options, ", ")
			}
			entry.writeField(output, options, key, field)
			index++
		}
		entry.writeString(output, options, ", ") // the Task ID follows
	}
	// Always write the Task ID at the end of he fields
	entry.writeString(output, options, "tid=")
	entry.writeInt64(output, options, entry.TaskID)
	entry.writeString(output, options, ")")

	log.Debugf("Blobs: %v", entry.Blobs)
	if len(entry.Blobs) > 0 {
		entry.writeString(output, options, "\n")
		for key, field := range entry.Blobs {
			entry.writeString(output, options, "    ")
			entry.writeString(output, options, key)
			entry.writeString(output, options, ": ")
			entry.writeBlob(output, options, field, 4)
			entry.writeString(output, options, "\n")
		}
	}
}

func (entry LogEntry) writeIndent(output io.Writer, _ *OutputOptions, indent int) {
	for i := 0; i < indent; i++ {
		_, _ = output.Write([]byte(" "))
	}
}

func (entry LogEntry) writeString(output io.Writer, _ *OutputOptions, value string) {
	_, _ = output.Write([]byte(value))
}

func (entry LogEntry) writeInt64(output io.Writer, _ *OutputOptions, value int64) {
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

func (entry LogEntry) writeBlob(output io.Writer, options *OutputOptions, blob any, indent int) {
	if value, ok := blob.(string); ok {
		entry.writeString(output, options, "\"")
		entry.writeString(output, options, value)
		entry.writeString(output, options, "\"")
	} else if value, ok := blob.(float64); ok {
		entry.writeString(output, options, strconv.FormatFloat(value, 'g', -1, 64))
	} else if value, ok := blob.(bool); ok {
		if value {
			entry.writeString(output, options, "true")
		} else {
			entry.writeString(output, options, "false")
		}
	} else if values, ok := blob.([]any); ok {
		entry.writeString(output, options, "[\n")
		entry.writeIndent(output, options, indent)
		for index, value := range values {
			if index > 0 {
				entry.writeString(output, options, ", \n")
				entry.writeIndent(output, options, indent)
			}
			entry.writeString(output, options, ", \n")
			entry.writeBlob(output, options, value, indent+2)
		}
		entry.writeString(output, options, "]")
	} else if values, ok := blob.(map[string]any); ok {
		entry.writeString(output, options, "{\n")
		index := 0
		for key, value := range values {
			entry.writeIndent(output, options, indent+2)
			entry.writeString(output, options, "\"")
			entry.writeString(output, options, key)
			entry.writeString(output, options, "\": ")
			entry.writeBlob(output, options, value, indent+2)
			if index < len(values)-1 {
				entry.writeString(output, options, ",\n")
			} else {
				entry.writeString(output, options, "\n")
			}
			index++
		}
		entry.writeIndent(output, options, indent-2)
		entry.writeString(output, options, "}")
	} else {
		entry.writeString(output, options, "!!!"+fmt.Sprintf("%v", blob))
	}
}

func (entry LogEntry) writeTimestamp(output io.Writer, options *OutputOptions) {
	timestamp := entry.Time.UTC()
	if options.Location != nil {
		timestamp = entry.Time.In(options.Location)
	}

	if options.Output.Value == "short" {
		timestampFormat := "15:04:05.000Z"
		if options.Location != nil {
			timestampFormat = "15:04:05.000"
		}
		entry.writeString(output, options, timestamp.Format(timestampFormat))
		entry.writeString(output, options, " ")
	} else {
		timestampFormat := "2006-01-02T15:04:05.000"
		if options.Location != nil {
			timestampFormat = "2006-01-02T15:04:05.000Z07:00"
		}
		entry.writeString(output, options, "[")
		entry.writeString(output, options, timestamp.Format(timestampFormat))
		entry.writeString(output, options, "] ")
	}
}

func (entry LogEntry) writeHeader(output io.Writer, options *OutputOptions) {
	entry.writeTimestamp(output, options)
	entry.Level.Write(output, options)

	if options.Output.Value == "short" {
		if len(entry.Name) > 0 {
			entry.writeString(output, options, " ")
			entry.writeString(output, options, entry.Name)
		}
	} else {
		entry.writeString(output, options, ": ")
		if len(entry.Name) > 0 {
			entry.writeString(output, options, entry.Name)
		}
		if entry.PID > 0 {
			entry.writeString(output, options, "/")
			entry.writeInt64(output, options, entry.PID)
		}
		if len(entry.Hostname) > 0 {
			entry.writeString(output, options, " on ")
			entry.writeString(output, options, entry.Hostname)
		}
	}
}

func (entry LogEntry) writeTopicAndScope(output io.Writer, options *OutputOptions) {
	if len(entry.Topic) > 0 {
		entry.writeStringWithColor(output, options, entry.Topic, Green)
		if len(entry.Scope) > 0 {
			entry.writeString(output, options, "/")
			entry.writeStringWithColor(output, options, entry.Scope, Yellow)
		}
		entry.writeString(output, options, " ")
	}
}

func (entry LogEntry) writeField(output io.Writer, options *OutputOptions, field string, value any) {
	entry.writeString(output, options, field)
	entry.writeString(output, options, "=")
	if value == nil {
		entry.writeString(output, options, "<null>")
	} else if actual, ok := value.(string); ok {
		entry.writeString(output, options, actual)
	} else if actual, ok := value.(float64); ok {
		entry.writeString(output, options, strconv.FormatFloat(actual, 'g', -1, 64))
	} else if actual, ok := value.(bool); ok {
		if actual {
			entry.writeString(output, options, "true")
		} else {
			entry.writeString(output, options, "false")
		}
	} else if values, ok := value.([]any); ok {
		entry.writeString(output, options, "[")
		for index, item := range values {
			if index > 0 {
				entry.writeString(output, options, ", ")
			}
			entry.writeField(output, options, "", item)
		}
		entry.writeString(output, options, "]")
	} else {
		entry.writeString(output, options, fmt.Sprintf("%v", value))
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
	entry.Blobs = map[string]any{}
	for key, value := range data {
		switch key {
		case "hostname":
			if entry.Hostname, ok = value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("hostname", value))
			}
		case "name":
			if entry.Name, ok = value.(string); !ok {
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
			if number, ok := value.(float64); ok {
				entry.Level = LogLevel(int(number))
			} else if str, ok := value.(string); ok {
				entry.Level = LogLevel(logger.ParseLevel(str))
			} else {
				merr.Append(errors.ArgumentInvalid.With("level", value))
			}
		case "pid":
			if number, ok := value.(float64); !ok {
				merr.Append(errors.ArgumentInvalid.With("pid", value))
			} else {
				entry.PID = int64(number)
			}
		case "tid":
			if number, ok := value.(float64); !ok {
				merr.Append(errors.ArgumentInvalid.With("tid", value))
			} else {
				entry.TaskID = int64(number)
			}
		case "time":
			if number, ok := value.(float64); ok {
				entry.Time = time.UnixMilli(int64(number))
			} else if tvalue, ok := value.(string); !ok {
				merr.Append(errors.ArgumentInvalid.With("time", value))
			} else if entry.Time, err = time.Parse(time.RFC3339, tvalue); err != nil {
				merr.Append(errors.Join(errors.ArgumentInvalid.With("time", value), err))
			}
		case "severity", "v":
			// ignore
		default:
			if value == nil {
				entry.Fields[key] = nil
			} else if _, ok := value.(string); ok {
				entry.Fields[key] = value
			} else if _, ok := value.(float64); ok {
				entry.Fields[key] = value
			} else if _, ok := value.(bool); ok {
				entry.Fields[key] = value
			} else if values, ok := value.([]any); ok && len(values) == 0 {
				entry.Fields[key] = values
			} else {
				entry.Blobs[key] = value
			}
		}
	}
	return merr.AsError()
}
