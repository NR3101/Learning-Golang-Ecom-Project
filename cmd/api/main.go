package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NR3101/go-ecom-project/internal/config"
	"github.com/NR3101/go-ecom-project/internal/database"
	"github.com/NR3101/go-ecom-project/internal/logger"
	"github.com/NR3101/go-ecom-project/internal/server"
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

	srv := server.New(cfg, db, logger)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a goroutine so that it doesn't block the graceful shutdown handling below.
	go func() {
		logger.Info().Msgf("Starting server on port:  %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exiting")
}
