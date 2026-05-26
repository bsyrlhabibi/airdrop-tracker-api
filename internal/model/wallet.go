package model

import "time"

type Wallet struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	AccountID uint      `json:"account_id" gorm:"index;not null"`
	Label     string    `json:"label" gorm:"not null"`
	Address   string    `json:"address" gorm:"not null"`
	Chain     string    `json:"chain" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Account *Account `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}
