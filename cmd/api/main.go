package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NR3101/go-ecom-project/internal/config"
	"github.com/NR3101/go-ecom-project/internal/database"
	"github.com/NR3101/go-ecom-project/internal/events"
	"github.com/NR3101/go-ecom-project/internal/interfaces"
	"github.com/NR3101/go-ecom-project/internal/logger"
	"github.com/NR3101/go-ecom-project/internal/providers"
	"github.com/NR3101/go-ecom-project/internal/repositories"
	"github.com/NR3101/go-ecom-project/internal/server"
	"github.com/NR3101/go-ecom-project/internal/services"
	"github.com/gin-gonic/gin"
)

// @title           E-commerce API
// @version         1.0
// @description     This is a sample server for an e-commerce application.
// @termsOfService  http://swaager.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  no-email@enail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description  Type "Bearer {token}" to authenticate requests

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database instance")
	}
	defer mainDB.Close()

	gin.SetMode(cfg.Server.GinMode)
	eventPublisher, err := events.NewEventPublisher(context.Background(), &cfg.Aws)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize event publisher")
	}
	defer eventPublisher.Close()

	userRepo := repositories.NewUserRepository(db)
	cartRepo := repositories.NewCartRepository(db)

	authService := services.NewAuthService(userRepo, cartRepo, cfg, eventPublisher)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db)
	cartService := services.NewCartService(db)
	orderService := services.NewOrderService(db)

	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider, err = providers.NewS3Provider(cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize S3 upload provider")
		}
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	uploadService := services.NewUploadService(uploadProvider)

	srv := server.New(
		cfg,
		&log,
		authService,
		productService,
		userService,
		uploadService,
		cartService,
		orderService)
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
		log.Info().Msgf("Starting server on port:  %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
		return
	}

	log.Info().Msg("Server exiting")
}
