package routes

import (
	"github.com/gin-gonic/gin"
)

func setupStaticRoutes(rg *gin.RouterGroup) {
	static := rg.Group("/static")
	{
		static.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "static ok"})
		})
	}
}
