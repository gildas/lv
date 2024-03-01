package cmd

import (
	"golang.org/x/term"
)

func rightpad(s string, length int) string {
	for len(s) < length {
		s = s + " "
	}
	return s
}

func leftpad(s string, length int) string {
	for len(s) < length {
		s = " " + s
	}
	return s
}

func isatty() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

