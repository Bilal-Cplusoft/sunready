package models

import (
	"time"
)

type User struct {
	ID                 int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"column:updated_at"`
	FirstName          string    `json:"first_name" gorm:"column:firstname"`
	LastName           string    `json:"last_name" gorm:"column:lastname"`
	Email              string    `json:"email" gorm:"column:email;uniqueIndex"`
	Password           string    `json:"-" gorm:"column:password"`
	PhoneNumber        string    `json:"phone_number" gorm:"column:phone_number"`
	Street             string    `json:"street" gorm:"column:street"`
	City               string    `json:"city" gorm:"column:city"`
	State              string    `json:"state" gorm:"column:state"`
	PostalCode         string    `json:"postal_code" gorm:"column:postal_code"`
	Country            string    `json:"country" gorm:"column:country"`
	UserType           UserType  `json:"user_type" gorm:"column:user_type"`
	HomeOwnershipType  string    `json:"home_ownership_type" gorm:"column:home_ownership_type" example:"owner"`
	AverageMonthlyBill float64   `json:"average_monthly_bill" gorm:"column:average_monthly_bill" example:"150.00"`
	UtilityProvider    string    `json:"utility_provider" gorm:"column:utility_provider" example:"PG&E"`
}


func (User) TableName() string {
	return "users"
}

type UserType int16

const (
	UserTypeAdmin    UserType = 0
	UserTypeCustomer UserType = 1
	UserTypeGeneral  UserType = 2
)
