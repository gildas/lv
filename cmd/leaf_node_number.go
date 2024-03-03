package cmd

import "strconv"

type NumberNode struct {
	Value float64
}

func (node NumberNode) GetValue(entry LogEntry) string {
	return strconv.FormatFloat(node.Value, 'g', -1, 64)
}

func (node NumberNode) String() string {
	return strconv.FormatFloat(node.Value, 'g', -1, 64)
}
