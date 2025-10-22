package models

import (
	"strings"
	"time"
)

type Customer struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	FirstName   string    `json:"first_name" gorm:"column:first_name;not null" example:"John"`
	LastName    string    `json:"last_name" gorm:"column:last_name;not null" example:"Smith"`
	Email       string    `json:"email" gorm:"column:email;uniqueIndex;not null" example:"john.smith@email.com"`
	PhoneNumber *string   `json:"phone_number" gorm:"column:phone_number" example:"+1-555-123-4567"`
	Address     string    `json:"address" gorm:"column:address;not null" example:"123 Main St, San Francisco, CA 94102"`
	City        *string   `json:"city" gorm:"column:city" example:"San Francisco"`
	State       *string   `json:"state" gorm:"column:state" example:"CA"`
	ZipCode     *string   `json:"zip_code" gorm:"column:zip_code" example:"94102"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active;default:true" example:"true"`
	PropertyType        *string  `json:"property_type" gorm:"column:property_type" example:"single_family"`
	RoofType           *string  `json:"roof_type" gorm:"column:roof_type" example:"asphalt_shingle"`
	HomeOwnershipType  *string  `json:"home_ownership_type" gorm:"column:home_ownership_type" example:"owner"`
	AverageMonthlyBill *float64 `json:"average_monthly_bill" gorm:"column:average_monthly_bill" example:"150.00"`
	UtilityProvider    *string  `json:"utility_provider" gorm:"column:utility_provider" example:"PG&E"`
	LeadSource         *string `json:"lead_source" gorm:"column:lead_source" example:"website"`
	ReferralCode       *string `json:"referral_code" gorm:"column:referral_code" example:"FRIEND2024"`
	Status             string    `json:"status" gorm:"column:status;default:'prospect'" example:"prospect"`
	Notes              *string   `json:"notes" gorm:"column:notes" example:"Interested in 10kW system"`
	PreferredContactMethod *string `json:"preferred_contact_method" gorm:"column:preferred_contact_method" example:"email"`
}

func (Customer) TableName() string {
	return "customers"
}

func (c *Customer) Sanitize() {
	c.FirstName = strings.TrimSpace(c.FirstName)
	c.LastName = strings.TrimSpace(c.LastName)
	c.Email = strings.TrimSpace(strings.ToLower(c.Email))
	c.Address = strings.TrimSpace(c.Address)
	if c.City != nil {
		*c.City = strings.TrimSpace(*c.City)
	}
	if c.State != nil {
		*c.State = strings.TrimSpace(strings.ToUpper(*c.State))
	}
}

func (c *Customer) Validate() error {
	if len(c.FirstName) == 0 || len(c.FirstName) > 100 {
		return ErrInvalidCustomerFirstName
	}
	if len(c.LastName) == 0 || len(c.LastName) > 100 {
		return ErrInvalidCustomerLastName
	}
	if len(c.Email) == 0 || len(c.Email) > 255 {
		return ErrInvalidCustomerEmail
	}
	if len(c.Address) == 0 || len(c.Address) > 500 {
		return ErrInvalidCustomerAddress
	}
	return nil
}

func (c *Customer) GetFullName() string {
	return c.FirstName + " " + c.LastName
}

// Customer status constants
const (
	CustomerStatusProspect   = "prospect"
	CustomerStatusQualified  = "qualified"
	CustomerStatusProposal   = "proposal"
	CustomerStatusContract   = "contract"
	CustomerStatusInstall    = "install"
	CustomerStatusComplete   = "complete"
	CustomerStatusCancelled  = "cancelled"
)

// Property type constants
const (
	PropertyTypeSingleFamily = "single_family"
	PropertyTypeMultiFamily  = "multi_family"
	PropertyTypeTownhouse    = "townhouse"
	PropertyTypeCondominium  = "condominium"
	PropertyTypeApartment    = "apartment"
	PropertyTypeCommercial   = "commercial"
)

// Roof type constants
const (
	RoofTypeAsphaltShingle = "asphalt_shingle"
	RoofTypeConcreteTile   = "concrete_tile"
	RoofTypeClayTile       = "clay_tile"
	RoofTypeMetal          = "metal"
	RoofTypeWoodShake      = "wood_shake"
	RoofTypeSlate          = "slate"
	RoofTypeFlat           = "flat"
)

// Home ownership type constants
const (
	HomeOwnershipOwner  = "owner"
	HomeOwnershipRenter = "renter"
	HomeOwnershipOther  = "other"
)

// Contact method constants
const (
	ContactMethodEmail = "email"
	ContactMethodPhone = "phone"
	ContactMethodText  = "text"
	ContactMethodAny   = "any"
)
