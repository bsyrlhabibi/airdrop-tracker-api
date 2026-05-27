package model

import "time"

// Task belongs to an AccountAirdrop (per-account task for an airdrop).
type Task struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	AccountAirdropID uint       `json:"account_airdrop_id" gorm:"index;not null"`
	CategoryID       *uint      `json:"category_id" gorm:"index"`
	Name             string     `json:"name" gorm:"not null"`
	Status           string     `json:"status" gorm:"default:pending"` // pending, ongoing, finish, edit
	Frequency        string     `json:"frequency" gorm:"default:once"` // once, daily, weekly, monthly
	Date             *time.Time `json:"date"`
	GasSpent         float64    `json:"gas_spent" gorm:"default:0"`
	TxHash           string     `json:"tx_hash"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relations
	AccountAirdrop *AccountAirdrop `json:"account_airdrop,omitempty" gorm:"foreignKey:AccountAirdropID"`
	Category       *Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}
