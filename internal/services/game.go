package services

import (
	"context"
	"errors"
	"gt/internal/repository"
	"gt/internal/security"
	"time"
)

type GameService struct {
	gameLoginRepo        *repository.GameLoginRepository
	gameLoginRequestRepo *repository.GameLoginRequestRepository
}

func NewGameService(gameLoginRepo *repository.GameLoginRepository, gameLoginRequestRepo *repository.GameLoginRequestRepository) *GameService {
	return &GameService{gameLoginRepo: gameLoginRepo, gameLoginRequestRepo: gameLoginRequestRepo}
}

var (
	ErrGameLoginRequestNotFound = errors.New("game login request not found")
	ErrGameLoginRequestUsed     = errors.New("game login request already used")
	ErrGameLoginCodeNotFound    = errors.New("game login code not found")
)

type CreatedGameLoginRequest struct {
	GameLoginRequest *repository.GameLoginRequest
	Token            string
}

func (s *GameService) CreateGameLoginRequest(ctx context.Context) (*CreatedGameLoginRequest, error) {
	token := security.GenerateToken()
	hashedToken, err := security.HashPassword(token)
	if err != nil {
		return nil, err
	}
	gameLoginRequest, err := s.gameLoginRequestRepo.Create(ctx, &repository.CreateGameLoginRequestRequest{
		Token: string(hashedToken),
	})
	if err != nil {
		return nil, err
	}
	return &CreatedGameLoginRequest{
		GameLoginRequest: gameLoginRequest,
		Token:            token,
	}, nil
}

func verifyGameLoginRequest(req *repository.GameLoginRequest) error {
	if req.User != nil || req.ExpiresAt.Before(time.Now()) {
		return ErrGameLoginRequestUsed
	}
	return nil
}

func (s *GameService) GetGameLoginRequest(ctx context.Context, id string) (*repository.GameLoginRequest, error) {
	req, err := s.gameLoginRequestRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrGameLoginRequestNotFound
	}
	if err := verifyGameLoginRequest(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *GameService) GetGameLoginRequestState(ctx context.Context, id string, token string) (*repository.GameLoginRequest, error) {
	req, err := s.gameLoginRequestRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req == nil || req.ExpiresAt.Before(time.Now()) || req.GameLogin != nil {
		return nil, ErrGameLoginRequestNotFound
	}
	if !security.CheckPasswordHash(token, req.Token) {
		return nil, ErrGameLoginRequestNotFound
	}
	return req, nil
}

func (s *GameService) Login(ctx context.Context, gameLoginRequestID string, user *repository.User) error {
	req, err := s.gameLoginRequestRepo.GetByID(ctx, gameLoginRequestID)
	if err != nil {
		return err
	}
	if req == nil {
		return ErrGameLoginRequestNotFound
	}
	if err := verifyGameLoginRequest(req); err != nil {
		return err
	}
	req.UserID = &user.ID
	return s.gameLoginRequestRepo.Update(ctx, req)
}

type GameLoginExchanged struct {
	GameLogin *repository.GameLogin
	Token     string
}

func (s *GameService) Exchange(ctx context.Context, gameRequestId string, token string) (*GameLoginExchanged, error) {
	req, err := s.gameLoginRequestRepo.GetByID(ctx, gameRequestId)
	if err != nil {
		return nil, err
	}
	if req == nil || req.ExpiresAt.Before(time.Now()) || req.GameLogin != nil {
		return nil, ErrGameLoginRequestNotFound
	}
	if !security.CheckPasswordHash(token, req.Token) {
		return nil, ErrGameLoginRequestNotFound
	}
	userID := req.UserID
	if userID == nil {
		return nil, errors.New("game login request has no user ID")
	}
	loginToken := security.GenerateToken()
	hashedToken, err := security.HashPassword(loginToken)
	if err != nil {
		return nil, err
	}
	gameLogin, err := s.gameLoginRepo.Create(ctx, &repository.CreateGameLoginRequest{
		UserID: *userID,
		Token:  string(hashedToken),
	})
	if err != nil {
		return nil, err
	}
	req.GameLogin = gameLogin
	if err := s.gameLoginRequestRepo.Update(ctx, req); err != nil {
		return nil, err
	}
	return &GameLoginExchanged{
		GameLogin: gameLogin,
		Token:     loginToken,
	}, nil
}
