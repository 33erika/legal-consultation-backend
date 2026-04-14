package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"legal-consultation/internal/config"
	"legal-consultation/internal/database"
	"legal-consultation/internal/handler"
	"legal-consultation/internal/middleware"
	"legal-consultation/internal/models"
	"legal-consultation/internal/repository"
	"legal-consultation/internal/service"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.Initialize(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化仓库
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	consultationRepo := repository.NewConsultationRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)

	// 初始化服务
	authSvc := service.NewAuthService(userRepo, &cfg.JWT)
	notificationSvc := service.NewNotificationService(&cfg.DingTalk)
	consultationSvc := service.NewConsultationService(consultationRepo, attachmentRepo, userRepo, notificationSvc)
	statisticsSvc := service.NewStatisticsService(consultationRepo)
	adminSvc := service.NewAdminService(userRepo, notificationSvc)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authSvc)
	consultationHandler := handler.NewConsultationHandler(consultationSvc)
	legalHandler := handler.NewLegalHandler(consultationSvc)
	statisticsHandler := handler.NewStatisticsHandler(statisticsSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)

	// 初始化 Gin
	r := gin.Default()

	// 中间件
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())

	// 路由
	setupRoutes(r, authHandler, consultationHandler, legalHandler, statisticsHandler, adminHandler, authSvc)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	consultationHandler *handler.ConsultationHandler,
	legalHandler *handler.LegalHandler,
	statisticsHandler *handler.StatisticsHandler,
	adminHandler *handler.AdminHandler,
	authSvc *service.AuthService,
) {
	// 健康检查
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	api := r.Group("/api/v1")
	{
		// 认证
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(authSvc), authHandler.GetCurrentUser)
		}

		// 需要认证的路由
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authSvc))
		{
			// 咨询
			consultations := protected.Group("/consultations")
			{
				consultations.POST("", consultationHandler.Create)
				consultations.GET("", consultationHandler.List)
				consultations.GET("/:id", consultationHandler.Get)
				consultations.POST("/:id/accept", middleware.RequireRole(models.RoleLegalStaff, models.RoleLegalHead), consultationHandler.Accept)
				consultations.POST("/:id/reply", middleware.RequireRole(models.RoleLegalStaff, models.RoleLegalHead), consultationHandler.Reply)
				consultations.POST("/:id/request-supplement", middleware.RequireRole(models.RoleLegalStaff, models.RoleLegalHead), consultationHandler.RequestSupplement)
				consultations.POST("/:id/close", middleware.RequireRole(models.RoleLegalStaff, models.RoleLegalHead), consultationHandler.Close)
				consultations.POST("/:id/transfer", middleware.RequireRole(models.RoleLegalStaff, models.RoleLegalHead), consultationHandler.Transfer)
				consultations.POST("/:id/rate", consultationHandler.Rate)
				consultations.GET("/:id/similar", consultationHandler.Similar)
				consultations.GET("/search", consultationHandler.Search)
			}

			// 法务工作台
			legal := protected.Group("/legal")
			legal.Use(middleware.RequireRole(models.RoleLegalStaff, models.RoleLegalHead))
			{
				legal.GET("/dashboard", legalHandler.Dashboard)
				legal.GET("/consultation-pool", legalHandler.ConsultationPool)
				legal.GET("/my-tasks", legalHandler.MyTasks)
				legal.GET("/staff-list", legalHandler.StaffList)
			}

			// 统计
			statistics := protected.Group("/statistics")
			statistics.Use(middleware.RequireRole(models.RoleLegalHead, models.RoleAdmin))
			{
				statistics.GET("/overview", statisticsHandler.Overview)
				statistics.GET("/category-distribution", statisticsHandler.CategoryDistribution)
				statistics.GET("/processing-efficiency", statisticsHandler.ProcessingEfficiency)
				statistics.GET("/export", statisticsHandler.Export)
				statistics.GET("/staff-workload", statisticsHandler.StaffWorkload)
			}

			// 案例库（待实现）
			protected.Group("/cases")

			// 管理
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole(models.RoleAdmin))
			{
				admin.GET("/users", adminHandler.ListUsers)
				admin.POST("/users", adminHandler.CreateUser)
				admin.PUT("/users/:id", adminHandler.UpdateUser)
				admin.POST("/users/:id/reset-password", adminHandler.ResetPassword)
				admin.PUT("/users/:id/toggle-status", adminHandler.ToggleUserStatus)

				admin.GET("/departments", adminHandler.ListDepartments)
				admin.POST("/departments", adminHandler.CreateDepartment)
				admin.PUT("/departments/:id", adminHandler.UpdateDepartment)
				admin.DELETE("/departments/:id", adminHandler.DeleteDepartment)

				admin.GET("/consultation-types", adminHandler.ListConsultationTypes)
				admin.PUT("/consultation-types/:type", adminHandler.UpdateConsultationType)

				admin.GET("/system-config", adminHandler.GetSystemConfig)
				admin.PUT("/system-config", adminHandler.UpdateSystemConfig)
				admin.POST("/system-config/test-notification", adminHandler.TestNotification)
			}
		}
	}
}
