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
		&model.AirdropTask{},
		&model.AccountAirdrop{},
		&model.Task{},
		&model.Wallet{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration done")

	// Fix legacy schema: remove account_id column from airdrops if it exists
	// SQLite doesn't reliably support DROP COLUMN on all versions
	// Use table recreation approach which works on ALL SQLite versions
	fixLegacyAirdropSchema()
}

func fixLegacyAirdropSchema() {
	// Check if account_id column exists
	var count int
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('airdrops') WHERE name = 'account_id'").Scan(&count)
	if count == 0 {
		return // already clean
	}

	log.Println("Fixing legacy airdrops schema: removing account_id column...")

	// Step 1: Create new table with correct schema
	DB.Exec(`CREATE TABLE IF NOT EXISTS airdrops_new (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		chain TEXT NOT NULL,
		category TEXT DEFAULT 'rumored',
		priority TEXT DEFAULT 'medium',
		status TEXT DEFAULT 'active',
		url TEXT,
		deadline DATETIME,
		notes TEXT,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	// Step 2: Copy data (only columns that exist in new schema)
	DB.Exec(`INSERT INTO airdrops_new (id, user_id, name, chain, category, priority, status, url, deadline, notes, created_at, updated_at)
		SELECT id, user_id, name, chain, category, priority, status, url, deadline, notes, created_at, updated_at FROM airdrops`)

	// Step 3: Drop old table and rename
	DB.Exec(`DROP TABLE IF EXISTS airdrops`)
	DB.Exec(`ALTER TABLE airdrops_new RENAME TO airdrops`)

	log.Println("Legacy airdrops schema fixed successfully")
}
