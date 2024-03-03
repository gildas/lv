package cmd

import "regexp"

type RegexNode struct {
	Regex *regexp.Regexp
}

func (node RegexNode) GetValue(entry LogEntry) string {
	return node.Regex.String()
}

func (node RegexNode) String() string {
	return node.Regex.String()
}
