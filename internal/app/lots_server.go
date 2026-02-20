package app

import (
	"lots-service/internal/config"
	"lots-service/internal/delivery/http_handlers"
	"lots-service/internal/repository"
	"lots-service/internal/server"
	"lots-service/internal/service"
	"lots-service/pkg/database"
)

func Run(cfg *config.Config) {
	db := database.NewPostgresConnection(cfg.DB.Host, cfg.DB.DBName, cfg.DB.User, cfg.DB.Password)

	repo := repository.NewPostgresLotsRepo(db)
	lotsService := service.NewLotsService(repo, cfg.StorageURL)
	lotsHandler := http_handlers.NewLotsHandler(lotsService)

	handler := server.NewRouter(lotsHandler)

	server.StartServer(handler, cfg.Port, cfg.Timeout)
}
