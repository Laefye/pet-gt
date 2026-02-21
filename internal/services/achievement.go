package services

import (
	"context"
	"errors"
	"gt/internal/repository"
)

type AchievementService struct {
	achievementRepo *repository.AchievementRepository
}

func (s *AchievementService) GetAchievementsByUserID(ctx context.Context, userID string) ([]*repository.Achievement, error) {
	return s.achievementRepo.GetByUserID(ctx, userID)
}

func NewAchievementService(achievementRepo *repository.AchievementRepository) *AchievementService {
	return &AchievementService{achievementRepo: achievementRepo}
}

var (
	ErrAchievementAlreadyExists = errors.New("achievement already exists for user")
)

func (s *AchievementService) CreateAchievement(ctx context.Context, req *repository.CreateAchievementRequest) (*repository.Achievement, error) {
	contains, err := s.achievementRepo.Contains(ctx, req.UserID, req.Name)
	if err != nil {
		return nil, err
	}
	if contains {
		return nil, ErrAchievementAlreadyExists
	}
	achievement, err := s.achievementRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return achievement, nil
}
