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

	// Drop old account_id column from airdrops (legacy schema fix)
	// SQLite 3.35+ supports ALTER TABLE DROP COLUMN
	if DB.Migrator().HasColumn(&model.Airdrop{}, "account_id") {
		log.Println("Dropping legacy account_id column from airdrops...")
		DB.Exec("ALTER TABLE airdrops DROP COLUMN account_id")
	}

	// Ensure default account exists for each user (auto-migration from old schema)
}
