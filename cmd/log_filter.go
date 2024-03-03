package cmd

type LogFilter interface {
	Filter(entry LogEntry) bool
}

type AllLogFilter struct{}

func (filter AllLogFilter) Filter(entry LogEntry) bool {
	return true
}
