package database

import (
	"log"

	"github.com/bsyrlhabibi/airdrop/internal/config"
	"github.com/bsyrlhabibi/airdrop/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	log.Println("Database connected")
}

func Migrate() {
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.Airdrop{},
		&model.AccountAirdrop{},
		&model.Task{},
		&model.Wallet{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration done")
}
