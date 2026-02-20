package server

import (
	"net/http"

	"github.com/NR3101/go-ecom-project/internal/config"
	"github.com/NR3101/go-ecom-project/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService
	userService    *services.UserService
	uploadService  *services.UploadService
	cartService    *services.CartService
}

func New(cfg *config.Config,
	db *gorm.DB,
	logger *zerolog.Logger,
	authService *services.AuthService,
	productService *services.ProductService,
	userService *services.UserService,
	uploadService *services.UploadService,
	cartService *services.CartService,
) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// API routes
	router.GET("/health", s.healthCheckHandler)
	router.Static("/uploads", "./uploads")

	// Group API routes under /api/v1
	api := router.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", s.getProfile)
				users.PUT("/profile", s.updateProfile)
			}

			// Category routes
			categories := protected.Group("/categories")
			{
				categories.POST("", s.adminMiddleware(), s.createCategory)
				categories.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categories.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			// Product routes
			products := protected.Group("/products")
			{
				products.POST("", s.adminMiddleware(), s.createProduct)
				products.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				products.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				products.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)
			}

			// Cart routes
			cart := protected.Group("/cart")
			{
				cart.GET("", s.getCart)
				cart.POST("/items", s.addToCart)
				cart.PUT("/items/:itemId", s.updateCartItem)
				cart.DELETE("/items/:itemId", s.removeCartItem)
			}
		}

		// Public routes
		api.GET("/categories", s.getCategories)
		api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProductByID)
	}

	return router
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// corsMiddleware is a middleware function that sets CORS headers to allow cross-origin requests.
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
