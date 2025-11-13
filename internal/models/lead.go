package models

import (
"time"
)


type LeadSource int

const (
LeadSourceLegacy  LeadSource = 0
LeadSourceDrone   LeadSource = 1
LeadSourceEarth   LeadSource = 2
LeadSourceFlyover LeadSource = 3
LeadSourceNone    LeadSource = 4
)


type LeadState int

const (
LeadStateProgress     LeadState = 0
LeadStateDone         LeadState = 1
LeadStateErrored      LeadState = 2
LeadStateInitialized  LeadState = 3
)


type Lead struct {
	ID        int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ProjectID int   `json:"project_id" gorm:"column:project_id;not null" example:"1"`
	UserID    *int  `json:"user_id" gorm:"column:user_id"`
	User      User  `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Latitude  float64 `json:"latitude" gorm:"column:latitude;not null" example:"37.7749"`
	Longitude float64 `json:"longitude" gorm:"column:longitude;not null" example:"-122.4194"`
	SystemSize   float64 `json:"system_size" gorm:"column:system_size" example:"10.5"`
	PanelCount   int     `json:"panel_count" gorm:"column:panel_count" example:"30"`
	HardwareType *string `json:"hardware_type" gorm:"column:hardware_type"`
	KwhUsage     float64 `json:"kwh_usage" gorm:"column:kwh_usage" example:"12000"`
	PanelId      int     `json:"panel_id" gorm:"column:panel_id" example:"1"`
	InverterId   int     `json:"inverter_id" gorm:"column:inverter_id" example:"1"`
	Consumption       []int   `json:"consumption" gorm:"-"`
	Period            string  `json:"period" gorm:"column:period"`
	TargetSolarOffset int     `json:"target_solar_offset" gorm:"column:target_solar_offset"`
	Mode              *string `json:"mode" gorm:"column:mode"`
	Unit              string  `json:"unit" gorm:"column:unit"`
	AnnualProduction float64 `json:"annual_production" gorm:"column:annual_production" example:"13000"`
    UtilityID     *int   `json:"utility_id" gorm:"column:utility_id" example:"1"`
    TariffID      *int   `json:"tariff_id" gorm:"column:tariff_id" example:"1"`
    ExternalID    *int   `json:"external_id" gorm:"column:external_id" example:"1"`
}

func (Lead) TableName() string {
return "leads"
}


func (l *Lead) Validate() error {
if l.Latitude < -90 || l.Latitude > 90 {
		return ErrInvalidLeadLatitude
}
if l.Longitude < -180 || l.Longitude > 180 {
		return ErrInvalidLeadLongitude
}
return nil
}
