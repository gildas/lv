package cmd

import (
	"github.com/gildas/go-errors"
)

type MatchNode struct {
	Left  LeafNode
	Right LeafNode
}

func CreateMatchNode(left, right LeafNode) (MatchNode, error) {
	if _, ok := right.(RegexNode); ok {
		return MatchNode{left, right}, nil
	}
	return MatchNode{}, errors.InvalidType.With("regular expression", right.String())
}

func (node MatchNode) Evaluate(entry LogEntry) bool {
	return node.Right.(RegexNode).Regex.MatchString(node.Left.GetValue(entry))
}
