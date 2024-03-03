package cmd

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gildas/go-logger"
	"golang.org/x/term"
)

func leftpad(s string, length int) string {
	var result strings.Builder
	var pad = length - len(s)
	for i := 0; i < pad; i++ {
		result.WriteString(" ")
	}
	result.WriteString(s)
	return result.String()
}

func isatty() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func GetPager(context context.Context) (output io.WriteCloser, close func(), err error) {
	log := logger.Must(logger.FromContext(context)).Child("pager", "get")
	location := os.Getenv("PAGER")
	if location == "NOPAGER" {
		return os.Stdout, func() {}, nil
	}
	if len(location) == 0 {
		if location, err = exec.LookPath("less"); err != nil {
			log.Warnf("Failed to find pager less: %s", err)
			location, err = exec.LookPath("more")
			if err != nil {
				log.Warnf("Failed to find pager more: %s", err)
				return os.Stdout, func() {}, nil
			}
		}
	}
	log.Infof("Using Pager %s", location)
	pager := exec.Command(location)
	output, err = pager.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to create pipe for pager: %s", err)
		return nil, nil, err
	}
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr
	if err := pager.Start(); err != nil {
		log.Fatalf("Failed to start pager %s", location, err)
		return nil, nil, err
	}
	log.Debugf("Pager started")
	close = func() {
		log.Debugf("Waiting for pager %s", location)
		output.Close()
		if err = pager.Wait(); err != nil {
			log.Fatalf("Failed to wait for pager %s", location, err)
			return
		}
		log.Infof("Pager done")
	}
	return output, close, nil
}
