package model

import "time"

type Account struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Color     string    `json:"color" gorm:"default:#3B82F6"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Wallets        []Wallet        `json:"wallets,omitempty" gorm:"foreignKey:AccountID"`
	AccountAirdrops []AccountAirdrop `json:"account_airdrops,omitempty" gorm:"foreignKey:AccountID"`
}
