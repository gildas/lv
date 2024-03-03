package cmd

import (
	"fmt"
	"regexp"
	"strconv"
)

type LeafNode interface {
	fmt.Stringer
	GetValue(entry LogEntry) string
}

func ParseLeafNode(value string) LeafNode {
	if len(value) > 0 {
		if value[0] == '.' {
			return FieldNode{Name: value[1:]}
		} else if value[0] == '"' && value[len(value)-1] == '"' {
			return ConstantNode{Value: value[1 : len(value)-1]}
		} else if value[0] == '/' && value[len(value)-1] == '/' {
			rex, err := regexp.Compile(value[1 : len(value)-1])
			if err == nil {
				return RegexNode{Regex: rex}
			}
		} else if value == "true" || value == "false" {
			return BooleanNode{Value: value == "true"}
		} else if number, err := strconv.ParseFloat(value, 64); err == nil {
			return NumberNode{Value: number}
		}
	}
	return ConstantNode{Value: value}
}
