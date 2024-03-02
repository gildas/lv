package cmd

type ConstantNode struct {
	Value string
}

func (node ConstantNode) GetValue(entry LogEntry) string {
	return node.Value
}

func (node ConstantNode) String() string {
	return node.Value
}
