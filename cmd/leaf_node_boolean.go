package cmd

type BooleanNode struct {
	Value bool
}

func (node BooleanNode) GetValue(entry LogEntry) string {
	if node.Value {
		return "true"
	}
	return "false"
}

func (node BooleanNode) String() string {
	if node.Value {
		return "true"
	}
	return "false"
}
