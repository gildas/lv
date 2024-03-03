package cmd

type ConditionLogFilter struct {
	Condition ConditionNode
}

func NewConditionFilter(condition string) (*ConditionLogFilter, error) {
	node, err := ParseCondition(condition)
	return &ConditionLogFilter{node}, err
}

func (filter ConditionLogFilter) Filter(entry LogEntry) bool {
	return filter.Condition.Evaluate(entry)
}
