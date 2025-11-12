package models

import (
	"time"
)


type Panel struct {
	ID           int      `json:"id" gorm:"primaryKey;column:id"`
	Manufacturer string   `json:"manufacturer" gorm:"column:manufacturer"`
	Model        string   `json:"model" gorm:"column:model"`
	Wattage      float64  `json:"wattage" gorm:"column:wattage"`
	LongSide     float64  `json:"longside" gorm:"column:longside"`
	ShortSide    float64  `json:"shortside" gorm:"column:shortside"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}
