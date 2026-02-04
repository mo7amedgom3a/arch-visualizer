package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/controllers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
)

func setupV1Routes(api *gin.RouterGroup, srv *server.Server) {
	v1 := api.Group("/v1")
	{
		// Static Routes
		staticCtrl := controllers.NewStaticController(srv.StaticDataService)
		staticGroup := v1.Group("/static")
		{
			staticGroup.GET("/providers", staticCtrl.ListProviders)
			staticGroup.GET("/resource-types", staticCtrl.ListResourceTypes)
			staticGroup.GET("/resource-models", staticCtrl.ListResourceModels)
			staticGroup.GET("/cloud-config", staticCtrl.ListCloudConfigs)
		}

		setupStaticRoutes(v1)

		// Controllers
		projectCtrl := controllers.NewProjectController(srv.ProjectService)
		userCtrl := controllers.NewUserController(srv.UserService)
		diagramCtrl := controllers.NewDiagramController(srv.PipelineOrchestrator, srv.DiagramService, srv.ArchitectureService, slog.Default())
		iamCtrl := controllers.NewIAMController(srv.IAMService)
		generationCtrl := controllers.NewGenerationController(srv.PipelineOrchestrator, slog.Default())

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
			projects.GET("", projectCtrl.ListProjects) // New List endpoint
			projects.GET("/:id", projectCtrl.GetProject)
			projects.PUT("/:id", projectCtrl.UpdateProject)
			projects.DELETE("/:id", projectCtrl.DeleteProject)
			projects.POST("/:id/duplicate", projectCtrl.DuplicateProject)
			projects.GET("/:id/versions", projectCtrl.GetProjectVersions)
			projects.POST("/:id/restore", projectCtrl.RestoreProjectVersion)

			// Architecture endpoints
			projects.GET("/:id/architecture", projectCtrl.GetArchitecture)
			projects.PUT("/:id/architecture", projectCtrl.UpdateArchitecture)
			projects.PATCH("/:id/architecture/nodes/:nodeId", projectCtrl.UpdateArchitectureNode)
			projects.DELETE("/:id/architecture/nodes/:nodeId", projectCtrl.DeleteArchitectureNode)
			projects.POST("/:id/architecture/validate", projectCtrl.ValidateArchitecture)

			// Code Generation endpoints
			projects.POST("/:id/generate", generationCtrl.GenerateCode)
			projects.GET("/:id/download", generationCtrl.DownloadCode)
		}

		// IAM Routes
		iam := v1.Group("/iam")
		{
			iam.GET("/policies", iamCtrl.ListPolicies)
			iam.POST("/users", iamCtrl.CreateUser)
			iam.POST("/roles", iamCtrl.CreateRole)
		}

		// Diagrams Routes
		diagrams := v1.Group("/diagrams")
		{
			diagrams.POST("/process", diagramCtrl.ProcessDiagram)
			diagrams.POST("/validate", diagramCtrl.ValidateDiagram)
			diagrams.POST("/validate-rules", diagramCtrl.ValidateDomainRules)
		}
	}
}
