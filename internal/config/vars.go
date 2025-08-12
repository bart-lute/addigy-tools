package config

import "log/slog"

var (
	// Version if you see this version number, you are not running a version created through the Makefile
	// Tip: make go-build
	Version  = "xx.xx.xx"
	LogLevel = "error"

	LogLevels = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)
