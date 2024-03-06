package cmd

import "context"

type MultiLogFilter struct {
	Filters []LogFilter
}

func (filter MultiLogFilter) IsEmpty() bool {
	return len(filter.Filters) == 0
}

func (filter MultiLogFilter) AsFilter() LogFilter {
	if len(filter.Filters) == 0 {
		return AllLogFilter{}
	}
	if len(filter.Filters) == 1 {
		return filter.Filters[0]
	}
	return filter
}

func (filter *MultiLogFilter) Add(filters ...LogFilter) *MultiLogFilter {
	filter.Filters = append(filter.Filters, filters...)
	return filter
}

func (filter MultiLogFilter) Filter(context context.Context, entry LogEntry) bool {
	for _, f := range filter.Filters {
		if !f.Filter(context, entry) {
			return false
		}
	}
	return true
}
