package model

import "time"

// AirdropTask is a task/checklist item that belongs to an Airdrop directly.
// This is the "what to do" for this airdrop — shared reference, not per-account.
type AirdropTask struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	AirdropID   uint       `json:"airdrop_id" gorm:"index;not null"`
	CategoryID  *uint      `json:"category_id" gorm:"index"`
	Name        string     `json:"name" gorm:"not null"`
	Status      string     `json:"status" gorm:"default:pending"` // pending, ongoing, finish, cancel, or custom
	Date        *time.Time `json:"date"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Airdrop   *Airdrop  `json:"airdrop,omitempty" gorm:"foreignKey:AirdropID"`
	Category  *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}
