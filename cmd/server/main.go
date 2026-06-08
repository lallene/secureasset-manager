package main

import (
	"net/http"
	"secureasset-manager/internal/change"
	"secureasset-manager/internal/cmdb"

	_ "secureasset-manager/docs"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/config"
	"secureasset-manager/internal/dashboard"
	"secureasset-manager/internal/database"
	"secureasset-manager/internal/incident"
	"secureasset-manager/internal/middleware"
	"secureasset-manager/internal/notification"
	"secureasset-manager/internal/seeder"
	"secureasset-manager/internal/service"
	"secureasset-manager/internal/site"
	ws "secureasset-manager/internal/websocket"
	"secureasset-manager/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.LoadConfig()
	logger.Init()
	database.Connect(cfg)

	if err := database.DB.AutoMigrate(
		&site.Site{},
		&service.Service{},
		&asset.Asset{},
		&incident.Incident{},
		&auth.User{},
		&incident.IncidentComment{},
		&incident.IncidentAttachment{},
		&notification.Notification{},
		&cmdb.Application{},
		&cmdb.Database{},
		&cmdb.ConfigurationRelation{},
		&change.ChangeRequest{},
	); err != nil {
		panic("Erreur migration : " + err.Error())
	}

	seeder.Seed()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)

	router.GET(
		"/users",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		auth.GetUsers,
	)

	router.POST(
		"/users",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		auth.CreateUser,
	)

	router.PUT(
		"/users/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		auth.UpdateUser,
	)

	router.DELETE(
		"/users/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		auth.DeleteUser,
	)

	router.POST(
		"/assets",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent"),
		asset.CreateAsset,
	)

	router.GET(
		"/assets",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		asset.GetAssets,
	)

	router.GET(
		"/assets/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		asset.GetAssetByID,
	)

	router.PUT(
		"/assets/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent"),
		asset.UpdateAsset,
	)

	router.DELETE(
		"/assets/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		asset.DeleteAsset,
	)

	router.POST(
		"/incidents",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Requester"),
		incident.CreateIncident,
	)

	router.GET(
		"/incidents",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		incident.GetIncidents,
	)

	router.GET(
		"/incidents/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		incident.GetIncidentByID,
	)

	router.PUT(
		"/incidents/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Requester"),
		incident.UpdateIncident,
	)

	router.DELETE(
		"/incidents/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Requester"),
		incident.DeleteIncident,
	)

	router.PUT(
		"/incidents/:id/take",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Agent"),
		incident.TakeIncident,
	)

	router.PUT(
		"/incidents/:id/resolve",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Agent"),
		incident.ResolveIncident,
	)

	router.PUT(
		"/incidents/:id/close",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Agent"),
		incident.CloseIncident,
	)

	router.GET(
		"/incidents/:id/comments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		incident.GetIncidentComments,
	)

	router.POST(
		"/incidents/:id/comments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Agent", "Requester"),
		incident.AddIncidentComment,
	)

	router.POST(
		"/incidents/:id/attachments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Agent", "Requester"),
		incident.UploadIncidentAttachment,
	)

	router.GET(
		"/incidents/:id/attachments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		incident.GetIncidentAttachments,
	)

	router.GET(
		"/attachments/:id/download",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		incident.DownloadIncidentAttachment,
	)

	router.DELETE(
		"/attachments/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		incident.DeleteIncidentAttachment,
	)

	router.GET(
		"/notifications",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		notification.GetNotifications,
	)

	router.PUT(
		"/notifications/:id/read",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		notification.MarkAsRead,
	)

	router.GET(
		"/dashboard/stats",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		dashboard.GetStats,
	)

	router.GET(
		"/sites",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		site.GetSites,
	)

	router.POST("/sites",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		site.CreateSite,
	)

	router.PUT("/sites/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		site.UpdateSite,
	)

	router.DELETE("/sites/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		site.DeleteSite,
	)

	router.GET(
		"/services",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		service.GetServices,
	)

	router.POST("/services",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		service.CreateService,
	)

	router.PUT("/services/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		service.UpdateService,
	)

	router.DELETE("/services/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		service.DeleteService,
	)

	router.GET(
		"/cmdb/applications",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		cmdb.GetApplications,
	)

	router.GET(
		"/cmdb/applications/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		cmdb.GetApplicationByID,
	)

	router.POST(
		"/cmdb/applications",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.CreateApplication,
	)

	router.PUT(
		"/cmdb/applications/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.UpdateApplication,
	)

	router.DELETE(
		"/cmdb/applications/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.DeleteApplication,
	)

	router.GET("/cmdb/databases",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		cmdb.GetDatabases,
	)

	router.POST("/cmdb/databases",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.CreateDatabase,
	)

	router.PUT("/cmdb/databases/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.UpdateDatabase,
	)

	router.DELETE("/cmdb/databases/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.DeleteDatabase,
	)

	router.GET("/cmdb/relations",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		cmdb.GetRelations,
	)

	router.POST("/cmdb/relations",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.CreateRelation,
	)

	router.PUT("/cmdb/relations/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.UpdateRelation,
	)

	router.DELETE("/cmdb/relations/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.DeleteRelation,
	)

	router.GET(
		"/cmdb/impact/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		cmdb.GetAssetImpact,
	)

	router.GET("/cmdb/business-services",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		cmdb.GetBusinessServices,
	)

	router.POST("/cmdb/business-services",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.CreateBusinessService,
	)

	router.PUT("/cmdb/business-services/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.UpdateBusinessService,
	)

	router.DELETE("/cmdb/business-services/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		cmdb.DeleteBusinessService,
	)
	router.GET(
		"/cmdb/graph",
		middleware.JWTAuthMiddleware(),
		cmdb.GetGraphHandler,
	)

	router.GET(
		"/changes",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		change.GetChanges,
	)

	router.POST(
		"/changes",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent"),
		change.CreateChange,
	)

	router.PUT(
		"/changes/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent"),
		change.UpdateChange,
	)

	router.DELETE(
		"/changes/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		change.DeleteChange,
	)

	router.GET(
		"/changes/:id/impact",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Agent", "Requester"),
		change.GetChangeImpact,
	)

	router.PUT("/changes/:id/submit", middleware.JWTAuthMiddleware(), middleware.RequireRole("Admin", "Agent"), change.SubmitChange)
	router.PUT("/changes/:id/approve", middleware.JWTAuthMiddleware(), middleware.RequireRole("Admin"), change.ApproveChange)
	router.PUT("/changes/:id/reject", middleware.JWTAuthMiddleware(), middleware.RequireRole("Admin"), change.RejectChange)
	router.PUT("/changes/:id/implement", middleware.JWTAuthMiddleware(), middleware.RequireRole("Admin", "Agent"), change.ImplementChange)
	router.PUT("/changes/:id/close", middleware.JWTAuthMiddleware(), middleware.RequireRole("Admin", "Agent"), change.CloseChange)

	router.GET("/ws", ws.HandleWS)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		panic("Erreur démarrage serveur : " + err.Error())
	}
}
