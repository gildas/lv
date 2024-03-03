package cmd

import "github.com/gildas/go-logger"

type LevelLogFilter struct {
	LevelSet logger.LevelSet
}

func NewLevelLogFilter(level string) *LevelLogFilter {
	return &LevelLogFilter{LevelSet: logger.ParseLevels(level)}
}

func (filter LevelLogFilter) Filter(entry LogEntry) bool {
	return filter.LevelSet.ShouldWrite(logger.Level(entry.Level), entry.Topic, entry.Scope)
}
