package cmd

import "strings"

type EqualsNode struct {
	Left  LeafNode
	Right LeafNode
}

func CreateEqualNode(left, right string) EqualsNode {
	var leftNode, rightNode LeafNode
	if strings.HasPrefix(left, ".") {
		leftNode = FieldNode{strings.TrimPrefix(left, ".")}
	} else {
		leftNode = ConstantNode{left}
	}
	if strings.HasPrefix(right, ".") {
		rightNode = FieldNode{strings.TrimPrefix(right, ".")}
	} else {
		rightNode = ConstantNode{right}
	}
	return EqualsNode{leftNode, rightNode}
}

func (node EqualsNode) Evaluate(entry LogEntry) bool {
	return node.Left.GetValue(entry) == node.Right.GetValue(entry)
}
