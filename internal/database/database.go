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
		&model.Category{},
		&model.Airdrop{},
		&model.AirdropTask{},
		&model.AccountAirdrop{},
		&model.Task{},
		&model.Wallet{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration done")

	// Fix legacy schema issues
	fixLegacyAirdropSchema()
	fixLegacyAirdropTaskSchema()
	fixLegacyTaskSchema()
}

func fixLegacyAirdropSchema() {
	var count int
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('airdrops') WHERE name = 'account_id'").Scan(&count)
	if count == 0 {
		return
	}

	log.Println("Fixing legacy airdrops schema: removing account_id column...")

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

	DB.Exec(`INSERT INTO airdrops_new (id, user_id, name, chain, category, priority, status, url, deadline, notes, created_at, updated_at)
		SELECT id, user_id, name, chain, category, priority, status, url, deadline, notes, created_at, updated_at FROM airdrops`)

	DB.Exec(`DROP TABLE IF EXISTS airdrops`)
	DB.Exec(`ALTER TABLE airdrops_new RENAME TO airdrops`)

	log.Println("Legacy airdrops schema fixed successfully")
}

func fixLegacyAirdropTaskSchema() {
	// Check if old columns exist (description, is_completed, completed_at, frequency)
	var hasDesc, hasCompleted int
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('airdrop_tasks') WHERE name = 'description'").Scan(&hasDesc)
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('airdrop_tasks') WHERE name = 'is_completed'").Scan(&hasCompleted)

	if hasDesc == 0 && hasCompleted == 0 {
		return // already new schema
	}

	log.Println("Fixing legacy airdrop_tasks schema...")

	DB.Exec(`CREATE TABLE IF NOT EXISTS airdrop_tasks_new (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		airdrop_id INTEGER NOT NULL,
		category_id INTEGER,
		name TEXT NOT NULL DEFAULT '',
		status TEXT DEFAULT 'pending',
		date DATETIME,
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	// Migrate data: description -> name, is_completed -> status
	DB.Exec(`INSERT INTO airdrop_tasks_new (id, airdrop_id, name, status, sort_order, created_at, updated_at)
		SELECT id, airdrop_id,
			COALESCE(description, ''),
			CASE WHEN is_completed = 1 THEN 'finish' ELSE 'pending' END,
			sort_order, created_at, updated_at
		FROM airdrop_tasks`)

	DB.Exec(`DROP TABLE IF EXISTS airdrop_tasks`)
	DB.Exec(`ALTER TABLE airdrop_tasks_new RENAME TO airdrop_tasks`)

	log.Println("Legacy airdrop_tasks schema fixed successfully")
}

func fixLegacyTaskSchema() {
	// Check if old columns exist (description, is_completed, completed_at, frequency)
	var hasDesc, hasCompleted int
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('tasks') WHERE name = 'description'").Scan(&hasDesc)
	DB.Raw("SELECT COUNT(*) FROM pragma_table_info('tasks') WHERE name = 'is_completed'").Scan(&hasCompleted)

	if hasDesc == 0 && hasCompleted == 0 {
		return // already new schema
	}

	log.Println("Fixing legacy tasks schema...")

	DB.Exec(`CREATE TABLE IF NOT EXISTS tasks_new (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_airdrop_id INTEGER NOT NULL,
		category_id INTEGER,
		name TEXT NOT NULL DEFAULT '',
		status TEXT DEFAULT 'pending',
		date DATETIME,
		gas_spent REAL DEFAULT 0,
		tx_hash TEXT,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	// Migrate data: description -> name, is_completed -> status
	DB.Exec(`INSERT INTO tasks_new (id, account_airdrop_id, name, status, gas_spent, tx_hash, created_at)
		SELECT id, account_airdrop_id,
			COALESCE(description, ''),
			CASE WHEN is_completed = 1 THEN 'finish' ELSE 'pending' END,
			gas_spent, tx_hash, created_at
		FROM tasks`)

	DB.Exec(`DROP TABLE IF EXISTS tasks`)
	DB.Exec(`ALTER TABLE tasks_new RENAME TO tasks`)

	log.Println("Legacy tasks schema fixed successfully")
}
