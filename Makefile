PROJECT_NAME := addigy-tools
PKG := github.com/bart-lute/$(PROJECT_NAME)
TAG := $(shell git describe --tags)
COMMON_LD_FLAGS := -s -w -X $(PKG)/internal/config.Version=$(TAG)

# Build locally (for testing, etc)
go-build:
	go build -o $(PROJECT_NAME) -ldflags "$(COMMON_LD_FLAGS) -X $(PKG)/internal/config.LogLevel=debug"
