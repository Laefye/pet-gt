package services

import (
	"context"
	"errors"
	"gt/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
}

type SignupRequest struct {
	Username string
	Password string
}

type SignupError struct {
	Message string
}

func (e *SignupError) Error() string {
	return e.Message
}

func (s *AuthService) Signup(ctx context.Context, req SignupRequest) (*repository.User, error) {
	existing, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &SignupError{Message: "Username already exists"}
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.userRepo.Create(ctx, req.Username, string(hashed))
}

type LoginRequest struct {
	Username  string
	Password  string
	UserAgent string
}

var ErrInvalidCredentials = errors.New("invalid username or password")

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*repository.Session, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return s.sessionRepo.Create(ctx, user.ID, req.UserAgent)
}

func (s *AuthService) Authenticate(ctx context.Context, sessionID string) (*repository.Session, error) {
	return s.sessionRepo.GetByID(ctx, sessionID)
}
