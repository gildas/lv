package cmd

import (
	"context"

	"github.com/gildas/go-logger"
)

type LevelLogFilter struct {
	LevelSet logger.LevelSet
}

func NewLevelLogFilter(level string) *LevelLogFilter {
	return &LevelLogFilter{LevelSet: logger.ParseLevels(level)}
}

func (filter LevelLogFilter) Filter(context context.Context, entry LogEntry) bool {
	log := logger.Must(logger.FromContext(context)).Child("filter", "filter", "type", "level")

	shouldWrite := filter.LevelSet.ShouldWrite(logger.Level(entry.Level), entry.Topic, entry.Scope)
	log.Debugf("Should %s be logged? %t", entry.Level, shouldWrite)
	return shouldWrite
}
