package model

import "time"

// AccountAirdrop links an Account to a global Airdrop.
// Each account can "select" airdrops and track its own progress.
type AccountAirdrop struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AccountID uint      `json:"account_id" gorm:"index;not null"`
	AirdropID uint      `json:"airdrop_id" gorm:"index;not null"`
	Status    string    `json:"status" gorm:"default:active"` // active, completed, paused
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Account *Account `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	Airdrop *Airdrop `json:"airdrop,omitempty" gorm:"foreignKey:AirdropID"`
	Tasks   []Task   `json:"tasks,omitempty" gorm:"foreignKey:AccountAirdropID"`
}
