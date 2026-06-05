package main

// @title SecureAsset Manager API
// @version 1.0
// @description API de gestion de parc informatique et incidents sécurité
// @host localhost:8080
// @BasePath /

import (
	"net/http"
	"secureasset-manager/internal/notification"
	"secureasset-manager/internal/seeder"

	_ "secureasset-manager/docs"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/config"
	"secureasset-manager/internal/dashboard"
	"secureasset-manager/internal/database"
	"secureasset-manager/internal/incident"
	"secureasset-manager/internal/middleware"
	"secureasset-manager/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)
import ws "secureasset-manager/internal/websocket"

func main() {

	cfg := config.LoadConfig()

	logger.Init()

	database.Connect(cfg)

	if err := database.DB.AutoMigrate(
		&asset.Asset{},
		&incident.Incident{},
		&auth.User{},
		&incident.IncidentComment{},
		&notification.Notification{},
		&incident.IncidentAttachment{},
	); err != nil {
		panic("Erreur migration : " + err.Error())
	}

	seeder.Seed()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.POST("/register", auth.Register)

	router.POST("/login", auth.Login)

	router.GET(
		"/users",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin"),
		auth.GetUsers,
	)

	router.GET(
		"/technicians",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		auth.GetTechnicians,
	)

	router.POST(
		"/assets",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician"),
		asset.CreateAsset,
	)

	router.GET(
		"/assets",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		asset.GetAssets,
	)

	router.GET(
		"/assets/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician"),
		asset.GetAssetByID,
	)

	router.PUT(
		"/assets/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician"),
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
		middleware.RequireRole("Viewer"),
		incident.CreateIncident,
	)

	router.GET(
		"/incidents",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		incident.GetIncidents,
	)

	router.GET(
		"/incidents/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		incident.GetIncidentByID,
	)

	router.PUT(
		"/incidents/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Viewer"),
		incident.UpdateIncident,
	)

	router.DELETE(
		"/incidents/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Viewer"),
		incident.DeleteIncident,
	)

	router.PUT(
		"/incidents/:id/take",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Technician"),
		incident.TakeIncident,
	)

	router.PUT(
		"/incidents/:id/resolve",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Technician"),
		incident.ResolveIncident,
	)

	router.PUT(
		"/incidents/:id/close",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Technician"),
		incident.CloseIncident,
	)

	router.GET(
		"/dashboard/stats",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		dashboard.GetStats,
	)

	router.POST(
		"/incidents/:id/react",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Viewer"),
		incident.ReactToIncident,
	)

	router.GET(
		"/incidents/:id/comments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		incident.GetIncidentComments,
	)

	router.POST(
		"/incidents/:id/comments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Technician", "Viewer"),
		incident.AddIncidentComment,
	)

	router.GET(
		"/notifications",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		notification.GetNotifications,
	)

	router.PUT(
		"/notifications/:id/read",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		notification.MarkAsRead,
	)

	router.POST(
		"/incidents/:id/attachments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Technician", "Viewer"),
		incident.UploadIncidentAttachment,
	)

	router.GET(
		"/incidents/:id/attachments",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		incident.GetIncidentAttachments,
	)

	router.GET(
		"/attachments/:id/download",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		incident.DownloadIncidentAttachment,
	)

	router.DELETE(
		"/attachments/:id",
		middleware.JWTAuthMiddleware(),
		middleware.RequireRole("Admin", "Technician", "Viewer"),
		incident.DeleteIncidentAttachment,
	)

	router.GET("/ws", ws.HandleWS)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		panic("Erreur démarrage serveur : " + err.Error())
	}
}
