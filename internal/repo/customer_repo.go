package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) Create(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *CustomerRepo) GetByID(ctx context.Context, id int) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepo) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepo) Update(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *CustomerRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Customer{}, id).Error
}

func (r *CustomerRepo) List(ctx context.Context, limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&customers).Error
	return customers, err
}

func (r *CustomerRepo) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&customers).Error
	return customers, err
}

func (r *CustomerRepo) Search(ctx context.Context, query string, limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	searchPattern := "%" + query + "%"
	err := r.db.WithContext(ctx).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR address ILIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&customers).Error
	return customers, err
}

func (r *CustomerRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Customer{}).Count(&count).Error
	return count, err
}

func (r *CustomerRepo) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Customer{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *CustomerRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	return r.db.WithContext(ctx).Model(&models.Customer{}).Where("id = ?", id).Update("status", status).Error
}