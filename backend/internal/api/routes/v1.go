package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/controllers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
)

func setupV1Routes(api *gin.RouterGroup, srv *server.Server) {
	v1 := api.Group("/v1")
	{
		// Static Routes
		staticCtrl := controllers.NewStaticController(srv.StaticDataService)
		v1.GET("/static/providers", staticCtrl.ListProviders)
		v1.GET("/static/resource-types", staticCtrl.ListResourceTypes)
		v1.GET("/static/resource-models", staticCtrl.ListResourceModels)

		setupStaticRoutes(v1)

		// Controllers
		projectCtrl := controllers.NewProjectController(srv.ProjectService)
		userCtrl := controllers.NewUserController(srv.UserService)
		diagramCtrl := controllers.NewDiagramController(srv.PipelineOrchestrator)

		// Users Routes
		users := v1.Group("/users")
		{
			users.POST("", userCtrl.CreateUser)
			users.GET("/:id", userCtrl.GetUser)
			users.GET("/:id/projects", projectCtrl.ListUserProjects)
		}

		// Projects Routes
		projects := v1.Group("/projects")
		{
			projects.POST("", projectCtrl.CreateProject)
			projects.GET("/:id", projectCtrl.GetProject)
			projects.PUT("/:id", projectCtrl.UpdateProject)
			projects.DELETE("/:id", projectCtrl.DeleteProject)
		}

		// Diagrams Routes
		diagrams := v1.Group("/diagrams")
		{
			diagrams.POST("/process", diagramCtrl.ProcessDiagram)
		}
	}
}
