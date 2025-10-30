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
ID                  int        `json:"id" gorm:"primaryKey;column:id"`
CreatedAt           time.Time  `json:"created_at" gorm:"column:created_at"`
UpdatedAt           time.Time  `json:"updated_at" gorm:"column:updated_at"`
State               int        `json:"state" gorm:"column:state;not null;default:0" example:"0"`
ProjectID          int        `json:"project_id" gorm:"column:project_id;not null" example:"1"`
CreatorID           *int       `json:"creator_id" gorm:"column:creator_id" example:"1"`
Latitude            float64    `json:"latitude" gorm:"column:latitude;not null" example:"37.7749"`
Longitude           float64    `json:"longitude" gorm:"column:longitude;not null" example:"-122.4194"`
Address             string     `json:"address" gorm:"column:address" example:"123 Solar St, San Francisco, CA 94102"`
Source              int        `json:"source" gorm:"column:source;not null;default:0" example:"0"`
KwhUsage            float64    `json:"kwh_usage" gorm:"column:kwh_usage" example:"12000"`
KwhPerKwManual      int        `json:"kwh_per_kw_manual" gorm:"column:kwh_per_kw_manual" example:"1200"`
ElectricityCostPre  *int       `json:"electricity_cost_pre" gorm:"column:electricity_cost_pre" example:"150"`
ElectricityCostPost *int       `json:"electricity_cost_post" gorm:"column:electricity_cost_post" example:"50"`
AdditionalIncentive *int       `json:"additional_incentive" gorm:"column:additional_incentive" example:"1000"`
SystemSize          float64    `json:"system_size" gorm:"column:system_size" example:"10.5"`
PanelCount          int        `json:"panel_count" gorm:"column:panel_count" example:"30"`
PanelID             *int       `json:"panel_id" gorm:"column:panel_id" example:"1"`
InverterID          *int       `json:"inverter_id" gorm:"column:inverter_id" example:"1"`
InverterCount       int        `json:"inverter_count" gorm:"column:inverter_count;default:1" example:"1"`
BatteryCount        int        `json:"battery_count" gorm:"column:battery_count;default:0" example:"0"`
UtilityID           *int       `json:"utility_id" gorm:"column:utility_id" example:"1"`
TariffID            *int       `json:"tariff_id" gorm:"column:tariff_id" example:"1"`
RoofMaterial        *int       `json:"roof_material" gorm:"column:roof_material" example:"1"`
SurfaceID           *int       `json:"surface_id" gorm:"column:surface_id" example:"1"`
AnnualProduction    float64    `json:"annual_production" gorm:"column:annual_production" example:"13000"`
WelcomeCallState       *int `json:"welcome_call_state" gorm:"column:welcome_call_state" example:"0"`
FinancingState         *int `json:"financing_state" gorm:"column:financing_state" example:"0"`
UtilityBillState       *int `json:"utility_bill_state" gorm:"column:utility_bill_state" example:"0"`
DesignApprovedState    *int `json:"design_approved_state" gorm:"column:design_approved_state" example:"0"`
PermittingApprovedState *int `json:"permitting_approved_state" gorm:"column:permitting_approved_state" example:"0"`
SitePhotosState        *int `json:"site_photos_state" gorm:"column:site_photos_state" example:"0"`
InstallCrewState       *int `json:"install_crew_state" gorm:"column:install_crew_state" example:"0"`
InstallationState      *int `json:"installation_state" gorm:"column:installation_state" example:"0"`
FinalInspectionState   *int `json:"final_inspection_state" gorm:"column:final_inspection_state" example:"0"`
PtoState               *int `json:"pto_state" gorm:"column:pto_state" example:"0"`
InstallationDate *string `json:"installation_date" gorm:"column:installation_date" example:"2025-10-15"`
DateNtp          *string `json:"date_ntp" gorm:"column:date_ntp" example:"2025-10-10"`
DateInstalled    *string `json:"date_installed" gorm:"column:date_installed" example:"2025-10-15"`
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
