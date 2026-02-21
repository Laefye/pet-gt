package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"gt/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type GameService struct {
	gameLoginRepo        *repository.GameLoginRepository
	gameLoginCodeRepo    *repository.GameLoginCodeRepository
	gameLoginRequestRepo *repository.GameLoginRequestRepository
}

func NewGameService(gameLoginRepo *repository.GameLoginRepository, gameLoginCodeRepo *repository.GameLoginCodeRepository, gameLoginRequestRepo *repository.GameLoginRequestRepository) *GameService {
	return &GameService{gameLoginRepo: gameLoginRepo, gameLoginCodeRepo: gameLoginCodeRepo, gameLoginRequestRepo: gameLoginRequestRepo}
}

var (
	ErrGameLoginRequestNotFound = errors.New("game login request not found")
	ErrGameLoginRequestUsed     = errors.New("game login request already used")
	ErrGameLoginCodeNotFound    = errors.New("game login code not found")
)

func createToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate token: " + err.Error())
	}
	return hex.EncodeToString(b)
}

type CreatedGameLoginRequest struct {
	GameLoginRequest *repository.GameLoginRequest
	Token            string
}

func (s *GameService) CreateGameLoginRequest(ctx context.Context) (*CreatedGameLoginRequest, error) {
	token := createToken()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
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
	if req.GameLoginCode != nil || req.ExpiresAt.Before(time.Now()) {
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
	if req == nil || req.ExpiresAt.Before(time.Now()) || req.GameLoginCode != nil {
		return nil, ErrGameLoginRequestNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(req.Token), []byte(token)); err != nil {
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
	gameLogin, err := s.gameLoginCodeRepo.Create(ctx, &repository.CreateGameLoginCodeRequest{
		UserID: user.ID,
	})
	req.GameLoginCode = gameLogin
	return s.gameLoginRequestRepo.Update(ctx, req)
}

type GameLoginExchanged struct {
	GameLogin *repository.GameLogin
	Token     string
}

func (s *GameService) Exchange(ctx context.Context, gameLoginCodeID string, userID string) (*GameLoginExchanged, error) {
	code, err := s.gameLoginCodeRepo.GetByIDAndUserID(ctx, gameLoginCodeID, userID)
	if err != nil {
		return nil, err
	}
	if code == nil || code.ExpiresAt.Before(time.Now()) || code.GameLoginID != nil {
		return nil, ErrGameLoginCodeNotFound
	}
	token := createToken()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	gameLogin, err := s.gameLoginRepo.Create(ctx, &repository.CreateGameLoginRequest{
		UserID: userID,
		Token:  string(hashedToken),
	})
	if err != nil {
		return nil, err
	}
	code.GameLoginID = &gameLogin.ID
	if err := s.gameLoginCodeRepo.Update(ctx, code); err != nil {
		return nil, err
	}
	return &GameLoginExchanged{
		GameLogin: gameLogin,
		Token:     token,
	}, nil
}
