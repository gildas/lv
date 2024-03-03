package cmd

type NotNode struct {
	Node ConditionNode
}

func (node NotNode) Evaluate(entry LogEntry) bool {
	return !node.Node.Evaluate(entry)
}
