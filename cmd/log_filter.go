package cmd

import "context"

type LogFilter interface {
	Filter(context context.Context, entry LogEntry) bool
}

type AllLogFilter struct{}

func (filter AllLogFilter) Filter(_ context.Context, _ LogEntry) bool {
	return true
}
