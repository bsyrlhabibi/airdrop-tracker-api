package model

import "time"

// AirdropTask is a task/checklist item that belongs to an Airdrop directly.
// This is the "what to do" for this airdrop — shared reference, not per-account.
type AirdropTask struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	AirdropID   uint       `json:"airdrop_id" gorm:"index;not null"`
	Description string     `json:"description" gorm:"not null"`
	Frequency   string     `json:"frequency" gorm:"default:once"` // once, daily, weekly, monthly
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Airdrop *Airdrop `json:"airdrop,omitempty" gorm:"foreignKey:AirdropID"`
}
