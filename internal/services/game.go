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
	gameLoginRepo *repository.GameLoginRequestRepository
}

func NewGameService(gameLoginRepo *repository.GameLoginRequestRepository) *GameService {
	return &GameService{gameLoginRepo: gameLoginRepo}
}

var (
	ErrGameLoginRequestNotFound = errors.New("game login request not found")
	ErrGameLoginRequestUsed     = errors.New("game login request already used")
	ErrGameLoginRequestExpired  = errors.New("game login request expired")
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
	gameLoginRequest, err := s.gameLoginRepo.Create(ctx, &repository.CreateGameLoginRequestRequest{
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
	if req.UserID != nil {
		return ErrGameLoginRequestUsed
	}
	if req.ExpiresAt.Before(time.Now()) {
		return ErrGameLoginRequestExpired
	}
	return nil
}

func (s *GameService) GetGameLoginRequest(ctx context.Context, id string) (*repository.GameLoginRequest, error) {
	req, err := s.gameLoginRepo.GetByID(ctx, id)
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
	req, err := s.gameLoginRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrGameLoginRequestNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(req.Token), []byte(token)); err != nil {
		return nil, ErrGameLoginRequestNotFound
	}
	return req, nil
}

func (s *GameService) Login(ctx context.Context, gameLoginRequestID string, user *repository.User) error {
	req, err := s.gameLoginRepo.GetByID(ctx, gameLoginRequestID)
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
	return s.gameLoginRepo.Update(ctx, req)
}
