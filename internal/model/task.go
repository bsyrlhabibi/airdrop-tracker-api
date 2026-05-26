package model

import "time"

// Task belongs to an AccountAirdrop (per-account task for an airdrop).
type Task struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	AccountAirdropID uint      `json:"account_airdrop_id" gorm:"index;not null"`
	Description      string    `json:"description" gorm:"not null"`
	Frequency        string    `json:"frequency" gorm:"default:once"` // once, daily, weekly, monthly
	IsCompleted      bool      `json:"is_completed" gorm:"default:false"`
	CompletedAt      *time.Time `json:"completed_at"`
	NextDue          *time.Time `json:"next_due"`
	GasSpent         float64   `json:"gas_spent" gorm:"default:0"`
	TxHash           string    `json:"tx_hash"`
	CreatedAt        time.Time `json:"created_at"`

	// Relations
	AccountAirdrop *AccountAirdrop `json:"account_airdrop,omitempty" gorm:"foreignKey:AccountAirdropID"`
}
