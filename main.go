package main

import (
	"context"
	"os"

	"github.com/gildas/go-logger"
	"github.com/gildas/lv/cmd"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if len(os.Getenv("LOG_DESTINATION")) == 0 {
		os.Setenv("LOG_DESTINATION", "nil")
	}
	log := logger.Create(APP, logger.EnvironmentPrefix("LV_"))
	defer log.Flush()
	cmd.RootCmd.Use = APP + " [falgs] <file>"
	cmd.RootCmd.Version = Version()
	if err := cmd.Execute(log.ToContext(context.Background())); err != nil {
		log.Fatalf("Failed to execute command", err)
		os.Exit(1)
	}
}
