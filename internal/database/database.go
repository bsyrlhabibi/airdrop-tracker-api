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
		&model.Task{},
		&model.Wallet{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration done")

	// Auto-create default accounts for existing users
	migrateExistingData()
}

// migrateExistingData creates a "Default" account for users who have
// airdrops/wallets without an account_id, and assigns existing data to it.
func migrateExistingData() {
	var users []model.User
	DB.Find(&users)

	for _, user := range users {
		// Check if user has any accounts
		var accountCount int64
		DB.Model(&model.Account{}).Where("user_id = ?", user.ID).Count(&accountCount)

		if accountCount == 0 {
			// Check if user has any airdrops or wallets
			var airdropCount int64
			DB.Model(&model.Airdrop{}).Where("user_id = ?", user.ID).Count(&airdropCount)

			var walletCount int64
			DB.Model(&model.Wallet{}).Where("user_id = ?", user.ID).Count(&walletCount)

			if airdropCount > 0 || walletCount > 0 {
				// Create default account
				account := model.Account{
					UserID: user.ID,
					Name:   "Default",
					Color:  "#3B82F6",
					Notes:  "Auto-migrated from existing data",
				}
				DB.Create(&account)

				// Assign existing airdrops to default account
				DB.Model(&model.Airdrop{}).
					Where("user_id = ? AND account_id = 0", user.ID).
					Update("account_id", account.ID)

				// Assign existing wallets to default account
				DB.Model(&model.Wallet{}).
					Where("user_id = ? AND account_id = 0", user.ID).
					Update("account_id", account.ID)

				log.Printf("Migrated existing data for user %d to Default account %d", user.ID, account.ID)
			}
		}
	}
}
