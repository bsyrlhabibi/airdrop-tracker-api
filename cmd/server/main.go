package main

import (
	"log"

	"github.com/bsyrlhabibi/airdrop/internal/config"
	"github.com/bsyrlhabibi/airdrop/internal/database"
	"github.com/bsyrlhabibi/airdrop/internal/router"

	_ "github.com/bsyrlhabibi/airdrop/docs" // swagger docs
)

// @title           Airdrop Tracker API
// @version         1.0
// @description     Personal airdrop task management API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Bita
// @contact.email  support@airdrop-tracker.local

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      airdrop-tracker-api.fly.dev
// @BasePath  /
// @schemes   https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}"
func main() {
	cfg := config.Load()

	database.Connect(cfg)
	database.Migrate()

	r := router.Setup(cfg)

	log.Printf("Server running on :%s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
