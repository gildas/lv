package cmd

type AndNode struct {
	Left  ConditionNode
	Right ConditionNode
}

func (node AndNode) Evaluate(entry LogEntry) bool {
	return node.Left.Evaluate(entry) && node.Right.Evaluate(entry)
}
