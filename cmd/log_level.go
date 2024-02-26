package cmd

import (
	"io"

	"github.com/gildas/go-logger"
)

type LogLevel int

func (level LogLevel) Write(output io.Writer, options *OutputOptions) {
	if options.UseColors {
		_, _ = output.Write([]byte(LevelColors[int(level)])) // Be sure to supprot levels not in the map
	}
	_, _ = output.Write([]byte(leftpad(logger.Level(int(level)).String(), 5)))
	if options.UseColors {
		_, _ = output.Write([]byte(Reset))
	}
}
