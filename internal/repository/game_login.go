package repository

import (
	"context"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GameLogin struct {
	ID     string `gorm:"primaryKey"`
	UserID string `gorm:"index,not null"`
	Token  string `gorm:"not null"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type GameLoginRepository struct {
	db *gorm.DB
}

func NewGameLoginRepository(db *gorm.DB) *GameLoginRepository {
	return &GameLoginRepository{db: db}
}

type CreateGameLoginRequest struct {
	UserID string
	Token  string
}

func (r *GameLoginRepository) Create(ctx context.Context, req *CreateGameLoginRequest) (*GameLogin, error) {
	gameLogin := &GameLogin{
		ID:     ulid.Make().String(),
		UserID: req.UserID,
		Token:  req.Token,
	}
	if err := r.db.WithContext(ctx).Create(gameLogin).Error; err != nil {
		return nil, err
	}
	return gameLogin, nil
}

func (r *GameLoginRepository) GetByID(ctx context.Context, id string) (*GameLogin, error) {
	var gameLogin GameLogin
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&gameLogin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &gameLogin, nil
}
