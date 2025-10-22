package service

import (
	"context"

	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
)

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) List(ctx context.Context, companyID int, limit, offset int) ([]*models.User, error) {
	return s.userRepo.List(ctx, limit, offset)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) FindByIDs(ctx context.Context, ids []int) ([]*models.User, error) {
	return s.userRepo.FindByIDs(ctx, ids)
}
