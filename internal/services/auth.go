package services

import (
	"context"
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

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type SignupError struct {
	Message string
}

func (e *SignupError) Error() string {
	return e.Message
}

func (s *AuthService) Signup(ctx context.Context, req SignupRequest) (*repository.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, &SignupError{Message: "Username already exists"}
	}
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	return s.userRepo.CreateUser(ctx, repository.CreateUserRequest{
		Username: req.Username,
		Password: hashedPassword,
	})
}

type LoginError struct {
	Message string
}

func (e *LoginError) Error() string {
	return e.Message
}

type LoginRequest struct {
	Username  string
	Password  string
	UserAgent string
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*repository.Session, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || !checkPasswordHash(req.Password, user.Password) {
		return nil, &LoginError{Message: "Invalid username or password"}
	}

	session, err := s.sessionRepo.CreateSession(ctx, repository.CreateSessionRequest{
		UserID:    user.ID,
		UserAgent: req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) Authenticate(ctx context.Context, sessionID string) (*repository.Session, error) {
	return s.sessionRepo.GetSessionByID(ctx, sessionID)
}
