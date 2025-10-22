package repo

import (
	"context"
	"fmt"

	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"gorm.io/gorm"
)


type LeadRepo struct {
	db *gorm.DB
}


func NewLeadRepo(db *gorm.DB) *LeadRepo {
	return &LeadRepo{db: db}
}


func (r *LeadRepo) Create(ctx context.Context, lead *models.Lead) error {
	if err := lead.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Create(lead)
	if result.Error != nil {
		return fmt.Errorf("failed to create lead: %w", result.Error)
	}

	return nil
}


func (r *LeadRepo) GetByID(ctx context.Context, id int) (*models.Lead, error) {
	var lead models.Lead
	result := r.db.WithContext(ctx).First(&lead, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrLeadNotFound
		}
		return nil, fmt.Errorf("failed to get lead: %w", result.Error)
	}

	return &lead, nil
}



func (r *LeadRepo) Update(ctx context.Context, lead *models.Lead) error {
	if err := lead.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Save(lead)
	if result.Error != nil {
		return fmt.Errorf("failed to update lead: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrLeadNotFound
	}

	return nil
}


func (r *LeadRepo) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Lead{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete lead: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrLeadNotFound
	}

	return nil
}


func (r *LeadRepo) List(ctx context.Context, customerID *int, creatorID *int, limit, offset int) ([]*models.Lead, int64, error) {
	var leads []*models.Lead
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Lead{})

	if customerID != nil {
		query = query.Where("customer_id = ?", *customerID)
	}
	if creatorID != nil {
		query = query.Where("creator_id = ?", *creatorID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count leads: %w", err)
	}

	result := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&leads)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list leads: %w", result.Error)
	}

	return leads, total, nil
}
