package tail

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/common"
)

// NewTailer returns a new tailer runner
func NewTailer() *common.Runner {
	return &common.Runner{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Exec: func(ctx context.Context, args []string, stdout, stderr io.Writer) error {
			log := logger.Must(logger.FromContext(ctx)).Child("tail", "exec")

			// TODO: We need a tail implementation that works on Windows, as the tail command is not available there. We could use the hpcloud/tail package, but it has some issues with file rotation and might not be suitable for all use cases. For now, we will just use the tail command on Unix-like systems.
			log.Infof("Executing: tail -F %v", args)
			cmd := exec.Command("tail", append([]string{"-F"}, args...)...)
			cmd.Env = os.Environ()
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			if err := cmd.Start(); err != nil {
				return err
			}
			return cmd.Wait()
		},
	}
}
