package model

import "time"

type Airdrop struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"index;not null"`
	Name      string     `json:"name" gorm:"not null"`
	Chain     string     `json:"chain" gorm:"not null"`
	Category  string     `json:"category" gorm:"default:rumored"`
	Priority  string     `json:"priority" gorm:"default:medium"`
	Status    string     `json:"status" gorm:"default:active"`
	URL       string     `json:"url"`
	Deadline  *time.Time `json:"deadline"`
	Notes     string     `json:"notes"`
	Tasks     []Task     `json:"tasks,omitempty" gorm:"foreignKey:AirdropID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
