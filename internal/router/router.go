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
	taskRepo := repository.NewTaskRepo(db)
	walletRepo := repository.NewWalletRepo(db)

	// Handlers
	authH := handler.NewAuthHandler(userRepo, cfg.JWTSecret)
	accountH := handler.NewAccountHandler(accountRepo)
	airdropH := handler.NewAirdropHandler(airdropRepo)
	taskH := handler.NewTaskHandler(taskRepo)
	walletH := handler.NewWalletHandler(walletRepo)
	dashboardH := handler.NewDashboardHandler(db)

	api := r.Group("/api")

	// Public
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)

	// Protected
	auth := api.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))

	// Accounts (Sybil)
	auth.GET("/accounts", accountH.List)
	auth.POST("/accounts", accountH.Create)
	auth.GET("/accounts/:id", accountH.Get)
	auth.PUT("/accounts/:id", accountH.Update)
	auth.DELETE("/accounts/:id", accountH.Delete)

	// Airdrops
	auth.GET("/airdrops", airdropH.List)
	auth.POST("/airdrops", airdropH.Create)
	auth.GET("/airdrops/:id", airdropH.Get)
	auth.PUT("/airdrops/:id", airdropH.Update)
	auth.DELETE("/airdrops/:id", airdropH.Delete)

	// Tasks
	auth.GET("/airdrops/:id/tasks", taskH.List)
	auth.POST("/airdrops/:id/tasks", taskH.Create)
	auth.PUT("/tasks/:id/complete", taskH.Complete)
	auth.PUT("/tasks/:id/reset", taskH.Reset)
	auth.DELETE("/tasks/:id", taskH.Delete)

	// Wallets
	auth.GET("/wallets", walletH.List)
	auth.POST("/wallets", walletH.Create)
	auth.DELETE("/wallets/:id", walletH.Delete)

	// Dashboard
	auth.GET("/dashboard", dashboardH.Summary)

	return r
}
