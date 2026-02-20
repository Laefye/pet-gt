package repository

import (
	"context"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GameLoginRequest struct {
	ID            string  `gorm:"primaryKey"`
	LoginedUserID *string `gorm:"index"`
	LoginedUser   *User   `gorm:"foreignKey:LoginedUserID;references:ID"`
}

type GameLoginRequestRepository struct {
	db *gorm.DB
}

func NewGameLoginRequestRepository(db *gorm.DB) *GameLoginRequestRepository {
	return &GameLoginRequestRepository{db: db}
}

func (r *GameLoginRequestRepository) CreateGameLoginRequest(ctx context.Context) (*GameLoginRequest, error) {
	gameLoginRequest := &GameLoginRequest{
		ID:            ulid.Make().String(),
		LoginedUserID: nil,
	}
	if err := r.db.WithContext(ctx).Create(gameLoginRequest).Error; err != nil {
		return nil, err
	}
	return gameLoginRequest, nil
}

func (r *GameLoginRequestRepository) GetGameLoginRequestByID(ctx context.Context, id string) (*GameLoginRequest, error) {
	var gameLoginRequest GameLoginRequest
	if err := r.db.WithContext(ctx).Preload("LoginedUser").Where(&GameLoginRequest{ID: id}).First(&gameLoginRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &gameLoginRequest, nil
}

func (r *GameLoginRequestRepository) UpdateGameLoginRequest(ctx context.Context, gameLoginRequest *GameLoginRequest) error {
	return r.db.WithContext(ctx).Save(gameLoginRequest).Error
}
