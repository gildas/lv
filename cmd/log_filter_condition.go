package cmd

import "context"

type ConditionLogFilter struct {
	Condition ConditionNode
}

func NewConditionFilter(condition string) (*ConditionLogFilter, error) {
	node, err := ParseCondition(condition)
	return &ConditionLogFilter{node}, err
}

func (filter ConditionLogFilter) Filter(context context.Context, entry LogEntry) bool {
	return filter.Condition.Evaluate(entry)
}
