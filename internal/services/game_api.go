package services

import (
	"context"
	"errors"
	"gt/internal/repository"
)

type GameAPIService struct {
	gameLoginRepo *repository.GameLoginRequestRepository
}

func NewGameAPIService(gameLoginRepo *repository.GameLoginRequestRepository) *GameAPIService {
	return &GameAPIService{gameLoginRepo: gameLoginRepo}
}

var (
	ErrGameLoginRequestNotFound = errors.New("game login request not found")
	ErrGameLoginRequestUsed     = errors.New("game login request already used")
)

func (s *GameAPIService) CreateGameLoginRequest(ctx context.Context) (*repository.GameLoginRequest, error) {
	return s.gameLoginRepo.Create(ctx)
}

func (s *GameAPIService) GetGameLoginRequestByID(ctx context.Context, id string) (*repository.GameLoginRequest, error) {
	req, err := s.gameLoginRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrGameLoginRequestNotFound
	}
	if req.AuthorizedUserID != nil {
		return nil, ErrGameLoginRequestUsed
	}
	return req, nil
}

func (s *GameAPIService) GetGameLoginRequestState(ctx context.Context, id string) (*repository.GameLoginRequest, error) {
	return s.gameLoginRepo.GetByID(ctx, id)
}

func (s *GameAPIService) Login(ctx context.Context, gameLoginRequestID string, user *repository.User) error {
	req, err := s.GetGameLoginRequestByID(ctx, gameLoginRequestID)
	if err != nil {
		return err
	}
	req.AuthorizedUserID = &user.ID
	return s.gameLoginRepo.Update(ctx, req)
}
