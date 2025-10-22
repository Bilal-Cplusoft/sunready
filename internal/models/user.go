package models

import (
	"time"
)

type User struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	FirstName   *string   `json:"first_name" gorm:"column:firstname"`
	LastName    *string   `json:"last_name" gorm:"column:lastname"`
	Email       string    `json:"email" gorm:"column:email;uniqueIndex"`
	Password    *string   `json:"-" gorm:"column:password"`
	PhoneNumber *string   `json:"phone_number" gorm:"column:phone_number"`
	Address     *string   `json:"address" gorm:"column:address"`
}

func (User) TableName() string {
	return "users"
}

type UserType int16

const (
	UserTypeUnknown   UserType = 0
	UserTypeAdmin     UserType = 1
	UserTypeSales     UserType = 2
	UserTypeSupport   UserType = 3
	UserTypeInstaller UserType = 4
	UserTypeCustomer  UserType = 5
)
