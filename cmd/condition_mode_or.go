package cmd

type OrNode struct {
	Left  ConditionNode
	Right ConditionNode
}

func (node OrNode) Evaluate(entry LogEntry) bool {
	return node.Left.Evaluate(entry) || node.Right.Evaluate(entry)
}
