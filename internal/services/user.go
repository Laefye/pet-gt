package services

import (
	"context"
	"gt/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*repository.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
