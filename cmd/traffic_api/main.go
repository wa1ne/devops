package main

import (
	"trafficlightAPI/internal/config"
	"trafficlightAPI/internal/handlers"
	"trafficlightAPI/internal/middleware/logger"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.InitLogger("", cfg.Env)
	handlers.Run(cfg, logger)
}
