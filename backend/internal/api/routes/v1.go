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

		// Cost Controller
		costCtrl := controllers.NewCostController(srv.PricingService, srv.ProjectService, srv.OptimizationService)

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
			projects.GET("", projectCtrl.ListProjects)
			projects.GET("/:id", projectCtrl.GetProject)
			projects.PUT("/:id", projectCtrl.UpdateProject)
			projects.DELETE("/:id", projectCtrl.DeleteProject)
			projects.POST("/:id/duplicate", projectCtrl.DuplicateProject)

			// Architecture (read-only snapshot lookup)
			projects.GET("/:id/architecture", projectCtrl.GetArchitecture)

			// ── Version CRUD ──────────────────────────────────────────────
			versions := projects.Group("/:id/versions")
			{
				versions.POST("", projectCtrl.CreateVersion)
				versions.GET("", projectCtrl.ListVersions)
				versions.GET("/latest", projectCtrl.GetLatestVersion)
				versions.GET("/:version_id", projectCtrl.GetVersionDetail)
				versions.GET("/:version_id/architecture", projectCtrl.GetVersionArchitecture)
				versions.DELETE("/:version_id", projectCtrl.DeleteVersion)

				// Version-scoped utility actions
				versions.POST("/:version_id/validate", projectCtrl.ValidateVersion)
				versions.POST("/:version_id/export/terraform", generationCtrl.GenerateCodeForVersion)
				versions.POST("/:version_id/estimate-cost", costCtrl.EstimateVersionCost)
			}

			// Code Generation (kept for non-version-scoped download convenience)
			projects.GET("/:id/download", generationCtrl.DownloadCode)
		}

		// IAM Routes
		iam := v1.Group("/iam")
		{
			iam.GET("/policies", iamCtrl.ListPolicies)
			iam.GET("/policies/between", iamCtrl.ListPoliciesBetweenServices)
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

		// AWS Resource Metadata Routes
		metadataCtrl := controllers.NewResourceMetadataController(srv.ResourceMetadataService)
		rulesCtrl := controllers.NewRulesController()

		// AWS Rules
		awsRules := v1.Group("/aws/rules")
		{
			awsRules.GET("", rulesCtrl.List)
			awsRules.GET("/:service", rulesCtrl.Get)
		}

		// Networking (backward-compatible)
		awsNetworking := v1.Group("/aws/networking")
		{
			awsNetworking.GET("/schemas", metadataCtrl.ListSchemas)
			awsNetworking.GET("/schemas/:resource", metadataCtrl.GetSchema)
		}

		// Compute
		awsCompute := v1.Group("/aws/compute")
		{
			awsCompute.GET("/schemas", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "compute"})
				metadataCtrl.ListSchemasByService(c)
			})
			awsCompute.GET("/schemas/:resource", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "compute"})
				metadataCtrl.GetSchemaByService(c)
			})
		}

		// Storage
		awsStorage := v1.Group("/aws/storage")
		{
			awsStorage.GET("/schemas", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "storage"})
				metadataCtrl.ListSchemasByService(c)
			})
			awsStorage.GET("/schemas/:resource", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "storage"})
				metadataCtrl.GetSchemaByService(c)
			})
		}

		// Database
		awsDatabase := v1.Group("/aws/database")
		{
			awsDatabase.GET("/schemas", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "database"})
				metadataCtrl.ListSchemasByService(c)
			})
			awsDatabase.GET("/schemas/:resource", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "database"})
				metadataCtrl.GetSchemaByService(c)
			})
		}

		// IAM
		awsIAM := v1.Group("/aws/iam")
		{
			awsIAM.GET("/schemas", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "iam"})
				metadataCtrl.ListSchemasByService(c)
			})
			awsIAM.GET("/schemas/:resource", func(c *gin.Context) {
				c.Params = append(c.Params, gin.Param{Key: "service", Value: "iam"})
				metadataCtrl.GetSchemaByService(c)
			})
		}
	}
}
