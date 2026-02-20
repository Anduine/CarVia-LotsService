package main

import (
	"log/slog"
	_ "lots-service/docs"
	"lots-service/internal/app"
	"lots-service/internal/config"
	"lots-service/internal/lib/logger"

	"os"
)

//	@title			Lots Service API
//	@version		1.0
//	@description	This is a lots microservice.

//	@host		localhost:3011
//	@BasePath	/

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	config := config.MustLoadConfig()

	logger.InitGlobalLogger(os.Stdout, slog.LevelDebug)

	app.Run(config)
}
