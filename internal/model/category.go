package model

import "time"

// Category is a task category (e.g. "Bridge", "Swap", "Staking", "Social", "Daily")
type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Color     string    `json:"color" gorm:"default:#6B7280"`
	CreatedAt time.Time `json:"created_at"`
}
