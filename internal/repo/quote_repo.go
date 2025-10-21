package repo

import (
	"gorm.io/gorm"
)

type QuoteRepo struct {
	db *gorm.DB
}

func NewQuoteRepo(db *gorm.DB) *QuoteRepo {
	return &QuoteRepo{db: db}
}
