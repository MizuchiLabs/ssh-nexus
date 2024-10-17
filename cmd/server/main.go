package main

import (
	"log/slog"
	"os"

	_ "github.com/MizuchiLabs/ssh-nexus/internal/migrations"
	"github.com/MizuchiLabs/ssh-nexus/internal/service"
	"github.com/lmittmann/tint"
)

// Set up global logger with specified configuration
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	if err := service.Server(); err != nil {
		slog.Error("Backend setup failed", "Error", err)
		return
	}
}
