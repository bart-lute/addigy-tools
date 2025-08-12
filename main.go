package main

import (
	"github.com/bart-lute/addigy-tools/internal/cmd"
	"github.com/bart-lute/addigy-tools/internal/config"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(config.LogLevels[config.LogLevel])
	cmd.Execute()
}
