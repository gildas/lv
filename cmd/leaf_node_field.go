package cmd

type FieldNode struct {
	Name string
}

func (node FieldNode) GetValue(entry LogEntry) string {
	return entry.GetField(node.Name)
}

func (node FieldNode) String() string {
	return node.Name
}
