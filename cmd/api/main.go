package main

import (
	"database/sql"

	"github.com/NR3101/go-ecom-project/internal/config"
	"github.com/NR3101/go-ecom-project/internal/database"
	"github.com/NR3101/go-ecom-project/internal/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := logger.New()
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get database instance")
	}
	defer func(mainDB *sql.DB) {
		err := mainDB.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to close database connection")
		}
	}(mainDB)

	gin.SetMode(cfg.Server.GinMode)
	logger.Info().Msg("Starting API server")
}
