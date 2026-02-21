package services

import (
	"context"
	"errors"
	"gt/internal/repository"
	"gt/internal/security"
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
	Email    string
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
	hashed, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	return s.userRepo.Create(ctx, &repository.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
	})
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
	if !security.CheckPasswordHash(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}
	return s.sessionRepo.Create(ctx, user.ID, req.UserAgent)
}

func (s *AuthService) Authenticate(ctx context.Context, sessionID string) (*repository.Session, error) {
	return s.sessionRepo.GetByID(ctx, sessionID)
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}
