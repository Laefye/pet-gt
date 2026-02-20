package services

import (
	"context"
	"gt/internal/repository"
)

type GameAPIService struct {
	gameLoginRepo *repository.GameLoginRequestRepository
}

func NewGameAPIService(gameLoginRepo *repository.GameLoginRequestRepository) *GameAPIService {
	return &GameAPIService{gameLoginRepo: gameLoginRepo}
}

func (s *GameAPIService) CreateGameLoginRequest(ctx context.Context) (*repository.GameLoginRequest, error) {
	return s.gameLoginRepo.CreateGameLoginRequest(ctx)
}

type GameAPIError struct {
	Message string
}

func (e *GameAPIError) Error() string {
	return e.Message
}

func (s *GameAPIService) GetGameLoginRequestByID(ctx context.Context, id string) (*repository.GameLoginRequest, error) {
	gameLoginRequest, err := s.gameLoginRepo.GetGameLoginRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if gameLoginRequest == nil {
		return nil, &GameAPIError{Message: "Game login request not found"}
	}
	if gameLoginRequest.LoginedUserID != nil {
		return nil, &GameAPIError{Message: "Game login request already used"}
	}
	return gameLoginRequest, nil
}

func (s *GameAPIService) GetGameLoginRequestState(ctx context.Context, id string) (*repository.GameLoginRequest, error) {
	return s.gameLoginRepo.GetGameLoginRequestByID(ctx, id)
}

func (s *GameAPIService) Login(ctx context.Context, gameLoginRequestID string, user *repository.User) error {
	gameLoginRequest, err := s.GetGameLoginRequestByID(ctx, gameLoginRequestID)
	if err != nil {
		return err
	}
	if gameLoginRequest.LoginedUserID != nil {
		return &GameAPIError{Message: "Game login request already used"}
	}
	gameLoginRequest.LoginedUserID = &user.ID
	return s.gameLoginRepo.UpdateGameLoginRequest(ctx, gameLoginRequest)
}
