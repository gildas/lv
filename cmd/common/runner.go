package common

import (
	"context"
	"io"
)

// Runner describes a command runner that can be used to execute commands and capture their output.
type Runner struct {
	Stdout io.Writer
	Stderr io.Writer
	Exec   func(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) error
}
