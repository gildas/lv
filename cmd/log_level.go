package cmd

import (
	"io"

	"github.com/gildas/go-logger"
)

type LogLevel logger.Level

func (level LogLevel) Write(output io.Writer, options *OutputOptions) {
	if options.UseColors {
		_, _ = output.Write([]byte(LevelColors[int(level)])) // Be sure to supprot levels not in the map
	}
	_, _ = output.Write([]byte(leftpad(logger.Level(level).String(), 5)))
	if options.UseColors {
		_, _ = output.Write([]byte(Reset))
	}
}

func (level LogLevel) String() string {
	return logger.Level(level).String()
}
