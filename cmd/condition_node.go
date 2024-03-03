package cmd

import (
	"strings"

	"github.com/gildas/go-errors"
)

type ConditionNode interface {
	Evaluate(entry LogEntry) bool
}

func ParseCondition(condition string) (ConditionNode, error) {
	condition = strings.TrimSpace(condition)
	if strings.HasPrefix(condition, "(") && strings.HasSuffix(condition, ")") {
		return ParseCondition(condition[1 : len(condition)-1])
	}

	if strings.HasPrefix(condition, "!") {
		node, err := ParseCondition(condition[1:])
		if err != nil {
			return nil, err
		}
		return NotNode{node}, nil
	}

	logical := map[string]func(left, right ConditionNode) ConditionNode{
		"&&": func(left, right ConditionNode) ConditionNode { return AndNode{left, right} },
		"||": func(left, right ConditionNode) ConditionNode { return OrNode{left, right} },
	}

	for operation, createNode := range logical {
		parts := strings.Split(condition, operation)
		if len(parts) > 1 {
			left, err := ParseCondition(parts[0])
			if err != nil {
				return nil, err
			}
			right, err := ParseCondition(parts[1])
			if err != nil {
				return nil, err
			}
			return createNode(left, right), nil
		}
	}

	compare := map[string]func(left, right LeafNode) (ConditionNode, error){
		"==": func(left, right LeafNode) (ConditionNode, error) { return EqualsNode{left, right}, nil },
		"=~": func(left, right LeafNode) (ConditionNode, error) { return CreateMatchNode(left, right) },
	}

	for operation, createNode := range compare {
		parts := strings.Split(condition, operation)
		if len(parts) > 1 {
			return createNode(ParseLeafNode(strings.TrimSpace(parts[0])), ParseLeafNode(strings.TrimSpace(parts[1])))
		}
	}
	return nil, errors.ArgumentInvalid.With("condition", condition)
}
