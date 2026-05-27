package router

import (
	"time"

	"github.com/bsyrlhabibi/airdrop/internal/config"
	"github.com/bsyrlhabibi/airdrop/internal/database"
	"github.com/bsyrlhabibi/airdrop/internal/handler"
	"github.com/bsyrlhabibi/airdrop/internal/middleware"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS — allow frontend origin
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	db := database.DB

	// Repositories
	userRepo := repository.NewUserRepo(db)
	accountRepo := repository.NewAccountRepo(db)
	airdropRepo := repository.NewAirdropRepo(db)
	aaRepo := repository.NewAccountAirdropRepo(db)
	taskRepo := repository.NewTaskRepo(db)
	walletRepo := repository.NewWalletRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)

	// Handlers
	authH := handler.NewAuthHandler(userRepo, cfg.JWTSecret)
	accountH := handler.NewAccountHandler(accountRepo, airdropRepo, aaRepo, db)
	airdropH := handler.NewAirdropHandler(airdropRepo)
	aaH := handler.NewAccountAirdropHandler(aaRepo)
	taskH := handler.NewTaskHandler(taskRepo, aaRepo)
	walletH := handler.NewWalletHandler(walletRepo)
	dashboardH := handler.NewDashboardHandler(db)
	categoryH := handler.NewCategoryHandler(categoryRepo)

	api := r.Group("/api")

	// Public
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)

	// Protected
	auth := api.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))

	// Categories
	auth.GET("/categories", categoryH.List)
	auth.POST("/categories", categoryH.Create)
	auth.PUT("/categories/:id", categoryH.Update)
	auth.DELETE("/categories/:id", categoryH.Delete)

	// Accounts (Sybil)
	auth.GET("/accounts", accountH.List)
	auth.POST("/accounts", accountH.Create)
	auth.GET("/accounts/:id", accountH.Get)
	auth.PUT("/accounts/:id", accountH.Update)
	auth.DELETE("/accounts/:id", accountH.Delete)

	// Account → Airdrop assignment
	auth.POST("/accounts/:id/airdrops", accountH.AssignAirdrop)
	auth.GET("/accounts/:id/airdrops", accountH.GetAccountAirdrops)
	auth.DELETE("/accounts/:id/airdrops/:airdrop_id", accountH.RemoveAirdrop)

	// Account clone
	auth.POST("/accounts/:id/clone", accountH.CloneAccount)

	// Airdrops (global catalog)
	auth.GET("/airdrops", airdropH.List)
	auth.POST("/airdrops", airdropH.Create)
	auth.GET("/airdrops/:id", airdropH.Get)
	auth.PUT("/airdrops/:id", airdropH.Update)
	auth.DELETE("/airdrops/:id", airdropH.Delete)

	// Airdrop Tasks (global tasks per airdrop)
	airdropTaskRepo := repository.NewAirdropTaskRepo(db)
	airdropTaskH := handler.NewAirdropTaskHandler(airdropTaskRepo, airdropRepo)
	auth.GET("/airdrops/:id/tasks", airdropTaskH.List)
	auth.POST("/airdrops/:id/tasks", airdropTaskH.Create)
	auth.POST("/airdrops/:id/tasks/bulk", airdropTaskH.BulkCreate)
	auth.PUT("/airdrops/:id/tasks/reorder", airdropTaskH.Reorder)
	auth.PUT("/airdrop-tasks/:id", airdropTaskH.Update)
	auth.DELETE("/airdrop-tasks/:id", airdropTaskH.Delete)

	// AccountAirdrops (direct operations)
	auth.GET("/account-airdrops/:id", aaH.Get)
	auth.PUT("/account-airdrops/:id", aaH.Update)

	// Tasks (per account-airdrop)
	auth.GET("/account-airdrops/:id/tasks", taskH.List)
	auth.POST("/account-airdrops/:id/tasks", taskH.Create)
	auth.PUT("/tasks/:id", taskH.Update)
	auth.DELETE("/tasks/:id", taskH.Delete)

	// Today's tasks (per account)
	auth.GET("/accounts/:id/tasks/today", taskH.TodayTasks)
	auth.GET("/accounts/:id/tasks/by-date", taskH.DateTasks)

	// Wallets
	auth.GET("/wallets", walletH.List)
	auth.POST("/wallets", walletH.Create)
	auth.DELETE("/wallets/:id", walletH.Delete)

	// Export
	exportH := handler.NewExportHandler(db)
	auth.GET("/export/excel", exportH.ExportExcel)

	// Dashboard
	auth.GET("/dashboard", dashboardH.Summary)
	auth.GET("/dashboard/comparison", accountH.GetComparison)

	return r
}
