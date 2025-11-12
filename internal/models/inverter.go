package models

import (
	"time"
)

type Inverter struct {
	ID           int       `json:"id" gorm:"primaryKey;column:id"`
	Manufacturer string    `json:"manufacturer" gorm:"column:manufacturer"`
	Model        string    `json:"model" gorm:"column:model"`
	Capacity     float64   `json:"capacity" gorm:"column:capacity"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}
