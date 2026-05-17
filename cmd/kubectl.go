package cmd

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/gildas/go-logger"
)

// NewKubectl
func NewKubectl() *Runner {
	return &Runner{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Exec: func(ctx context.Context, args []string, stdout, stderr io.Writer) error {
			log := logger.Must(logger.FromContext(ctx)).Child("kubectl", "exec")

			log.Infof("Executing: kubectl %v", args)
			cmd := exec.Command("kubectl", args...)
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
