package model

import "time"

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	AirdropID   uint       `json:"airdrop_id" gorm:"index;not null"`
	Description string     `json:"description" gorm:"not null"`
	Frequency   string     `json:"frequency" gorm:"default:once"`
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
	NextDue     *time.Time `json:"next_due"`
	WalletID    *uint      `json:"wallet_id"`
	GasSpent    float64    `json:"gas_spent" gorm:"default:0"`
	TxHash      string     `json:"tx_hash"`
	CreatedAt   time.Time  `json:"created_at"`
}
