package cmd

import (
	"context"
	"io"
)

type Runner struct {
	Stdout io.Writer
	Stderr io.Writer
	Exec   func(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) error
}
