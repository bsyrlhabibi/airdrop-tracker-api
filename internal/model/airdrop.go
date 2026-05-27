package model

import "time"

// Airdrop is a global airdrop opportunity (catalog entry).
// No account_id — this is shared across all accounts.
type Airdrop struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"index;not null"`
	Name      string     `json:"name" gorm:"not null"`
	Chain     string     `json:"chain" gorm:"not null"`
	Category  string     `json:"category" gorm:"default:rumored"`
	Priority  string     `json:"priority" gorm:"default:medium"`
	Status    string     `json:"status" gorm:"default:active"`
	URL       string     `json:"url"`
	DateStart *time.Time `json:"date_start"`
	DateEnd   *time.Time `json:"date_end"`
	Deadline  *time.Time `json:"deadline"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
