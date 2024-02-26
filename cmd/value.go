package cmd

import (
	"fmt"
	"io"
	"strconv"
)

func printIndent(indent int, output io.Writer) {
	for i := 0; i < indent; i++ {
		_, _ = output.Write([]byte(" "))
	}
}

func printField(field any, output io.Writer, indent int) {
	if value, ok := field.(string); ok {
		_, _ = output.Write([]byte(value))
	} else if value, ok := field.(int64); ok {
		_, _ = output.Write([]byte(strconv.FormatInt(value, 10)))
	} else if values, ok := field.([]any); ok {
		_, _ = output.Write([]byte("[\n"))
		printIndent(indent, output)
		for index, value := range values {
			if index > 0 {
				_, _ = output.Write([]byte(", \n"))
				printIndent(indent, output)
			}
			printField(value, output, indent+2)
		}
		_, _ = output.Write([]byte("]"))
	} else if values, ok := field.(map[string]any); ok {
		_, _ = output.Write([]byte("{\n"))
		printIndent(indent, output)
		index := 0
		for key, value := range values {
			_, _ = output.Write([]byte("\""))
			_, _ = output.Write([]byte(key))
			_, _ = output.Write([]byte("\": "))
			printField(value, output, indent+2)
			if index < len(values)-1 {
				_, _ = output.Write([]byte(",\n"))
				printIndent(indent, output)
			} else {
				_, _ = output.Write([]byte("\n"))
			}
			index++
		}
		printIndent(indent-2, output)
		_, _ = output.Write([]byte("}"))
	} else {
		_, _ = output.Write([]byte("!!!" + fmt.Sprintf("%v", field)))
	}
}
