package service

import (
	"context"
	"strings"
	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
	"slices"
)

type CustomerService struct {
	customerRepo *repo.CustomerRepo
}

func NewCustomerService(customerRepo *repo.CustomerRepo) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	customer.Sanitize()
	if err := customer.Validate(); err != nil {
		return err
	}
	if customer.Status == "" {
		customer.Status = models.CustomerStatusProspect
	}

	return s.customerRepo.Create(ctx, customer)
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, id int) (*models.Customer, error) {
	return s.customerRepo.GetByID(ctx, id)
}

func (s *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*models.Customer, error) {
	return s.customerRepo.GetByEmail(ctx, email)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customer *models.Customer) error {
	customer.Sanitize()
	if err := customer.Validate(); err != nil {
		return err
	}
	return s.customerRepo.Update(ctx, customer)
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id int) error {
	return s.customerRepo.Delete(ctx, id)
}

func (s *CustomerService) ListCustomers(ctx context.Context, limit, offset int) ([]*models.Customer, error) {
	return s.customerRepo.List(ctx, limit, offset)
}

func (s *CustomerService) ListCustomersByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Customer, error) {
	return s.customerRepo.ListByStatus(ctx, status, limit, offset)
}

func (s *CustomerService) SearchCustomers(ctx context.Context, query string, limit, offset int) ([]*models.Customer, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return s.customerRepo.List(ctx, limit, offset)
	}
	return s.customerRepo.Search(ctx, query, limit, offset)
}

func (s *CustomerService) GetCustomerCount(ctx context.Context) (int64, error) {
	return s.customerRepo.Count(ctx)
}

func (s *CustomerService) GetCustomerCountByStatus(ctx context.Context, status string) (int64, error) {
	return s.customerRepo.CountByStatus(ctx, status)
}

func (s *CustomerService) UpdateCustomerStatus(ctx context.Context, id int, status string) error {
	validStatuses := []string{
		models.CustomerStatusProspect,
		models.CustomerStatusQualified,
		models.CustomerStatusProposal,
		models.CustomerStatusContract,
		models.CustomerStatusInstall,
		models.CustomerStatusComplete,
		models.CustomerStatusCancelled,
	}

	isValid := slices.Contains(validStatuses, status)
	if !isValid {
		return models.ErrInvalidCustomerStatus
	}

	return s.customerRepo.UpdateStatus(ctx, id, status)
}

func (s *CustomerService) GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Customer, error) {
	return s.customerRepo.GetByPhoneNumber(ctx, phoneNumber)
}

func (s *CustomerService) GetCustomerStats(ctx context.Context) (*CustomerStats, error) {
	total, err := s.customerRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	prospects, err := s.customerRepo.CountByStatus(ctx, models.CustomerStatusProspect)
	if err != nil {
		return nil, err
	}

	qualified, err := s.customerRepo.CountByStatus(ctx, models.CustomerStatusQualified)
	if err != nil {
		return nil, err
	}

	proposals, err := s.customerRepo.CountByStatus(ctx, models.CustomerStatusProposal)
	if err != nil {
		return nil, err
	}

	contracts, err := s.customerRepo.CountByStatus(ctx, models.CustomerStatusContract)
	if err != nil {
		return nil, err
	}

	installations, err := s.customerRepo.CountByStatus(ctx, models.CustomerStatusInstall)
	if err != nil {
		return nil, err
	}

	completed, err := s.customerRepo.CountByStatus(ctx, models.CustomerStatusComplete)
	if err != nil {
		return nil, err
	}

	cancelled, err := s.customerRepo.CountByStatus(ctx, models.CustomerStatusCancelled)
	if err != nil {
		return nil, err
	}

	return &CustomerStats{
		Total:         total,
		Prospects:     prospects,
		Qualified:     qualified,
		Proposals:     proposals,
		Contracts:     contracts,
		Installations: installations,
		Completed:     completed,
		Cancelled:     cancelled,
	}, nil
}

type CustomerStats struct {
	Total         int64 `json:"total"`
	Prospects     int64 `json:"prospects"`
	Qualified     int64 `json:"qualified"`
	Proposals     int64 `json:"proposals"`
	Contracts     int64 `json:"contracts"`
	Installations int64 `json:"installations"`
	Completed     int64 `json:"completed"`
	Cancelled     int64 `json:"cancelled"`
}
