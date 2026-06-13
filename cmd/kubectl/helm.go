package kubectl

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd/common"
)

// NewHelm returns a new helm runner
func NewHelm() *common.Runner {
	return &common.Runner{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Exec: func(ctx context.Context, args []string, stdout, stderr io.Writer) error {
			log := logger.Must(logger.FromContext(ctx)).Child("helm", "exec")

			log.Infof("Executing: helm %v", args)
			cmd := exec.Command("helm", args...)
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

// IsHelmAvailable checks if helm is available in the system
func IsHelmAvailable() bool {
	_, err := exec.LookPath("helm")
	return err == nil
}
