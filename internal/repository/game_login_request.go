package repository

import (
	"context"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GameLoginRequest struct {
	ID              string         `gorm:"primaryKey"`
	Token           string         `gorm:"uniqueIndex"`
	GameLoginCodeID *string        `gorm:"index"`
	ExpiresAt       time.Time      `gorm:"not null"`
	GameLoginCode   *GameLoginCode `gorm:"foreignKey:GameLoginCodeID"`
}

type GameLoginRequestRepository struct {
	db *gorm.DB
}

func NewGameLoginRequestRepository(db *gorm.DB) *GameLoginRequestRepository {
	return &GameLoginRequestRepository{db: db}
}

type CreateGameLoginRequestRequest struct {
	Token string
}

func (r *GameLoginRequestRepository) Create(ctx context.Context, req *CreateGameLoginRequestRequest) (*GameLoginRequest, error) {
	gameLoginRequest := &GameLoginRequest{
		ID:        ulid.Make().String(),
		Token:     req.Token,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := r.db.WithContext(ctx).Create(gameLoginRequest).Error; err != nil {
		return nil, err
	}
	return gameLoginRequest, nil
}

func (r *GameLoginRequestRepository) GetByID(ctx context.Context, id string) (*GameLoginRequest, error) {
	var req GameLoginRequest
	err := r.db.WithContext(ctx).Preload("GameLoginCode").Preload("GameLoginCode.User").Where("id = ?", id).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *GameLoginRequestRepository) Update(ctx context.Context, req *GameLoginRequest) error {
	return r.db.WithContext(ctx).Save(req).Error
}
