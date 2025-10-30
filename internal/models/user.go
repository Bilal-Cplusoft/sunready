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
	UserType    UserType  `json:"user_type" gorm:"column:user_type"`
}

func (User) TableName() string {
	return "users"
}

type UserType int16

const (
	UserTypeAdmin     UserType = 0
	UserTypeCustomer  UserType = 1
	UserTypeGeneral   UserType = 2
)
