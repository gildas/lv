package cmd

import (
	"fmt"
	"regexp"
)

type LeafNode interface {
	fmt.Stringer
	GetValue(entry LogEntry) string
}

func ParseLeafNode(value string) LeafNode {
	if len(value) > 0 {
		if value[0] == '.' {
			return FieldNode{Name: value[1:]}
		} else if value[0] == '/' && value[len(value)-1] == '/' {
			rex, err := regexp.Compile(value[1 : len(value)-1])
			if err == nil {
				return RegexNode{Regex: rex}
			}
		}
	}
	return ConstantNode{Value: value}
}
