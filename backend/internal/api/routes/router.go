package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/docs" // Swagger Docs
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/middleware"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes the Gin router with middleware and routes
func SetupRouter(srv *server.Server) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger 
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Group
	api := r.Group("/api")
	{
		setupV1Routes(api, srv)
	}

	return r
}
