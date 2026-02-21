package repository

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GameLoginCode struct {
	ID          string     `gorm:"primaryKey"`
	UserID      string     `gorm:"index,not null"`
	GameLoginID *string    `gorm:"index"`
	ExpiresAt   time.Time  `gorm:"not null"`
	User        *User      `gorm:"foreignKey:UserID"`
	GameLogin   *GameLogin `gorm:"foreignKey:GameLoginID"`
}

type GameLoginCodeRepository struct {
	db *gorm.DB
}

func NewGameLoginCodeRepository(db *gorm.DB) *GameLoginCodeRepository {
	return &GameLoginCodeRepository{db: db}
}

type CreateGameLoginCodeRequest struct {
	UserID string
}

func (r *GameLoginCodeRepository) Create(ctx context.Context, req *CreateGameLoginCodeRequest) (*GameLoginCode, error) {
	gameLoginCode := &GameLoginCode{
		ID:        ulid.Make().String(),
		UserID:    req.UserID,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := r.db.WithContext(ctx).Preload("User").Create(gameLoginCode).Error; err != nil {
		return nil, err
	}
	return gameLoginCode, nil
}

func (r *GameLoginCodeRepository) GetByIDAndUserID(ctx context.Context, id string, userID string) (*GameLoginCode, error) {
	gameLoginCode := &GameLoginCode{}
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ? AND user_id = ?", id, userID).First(gameLoginCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return gameLoginCode, nil
}

func (r *GameLoginCodeRepository) Update(ctx context.Context, gameLoginCode *GameLoginCode) error {
	return r.db.WithContext(ctx).Save(gameLoginCode).Error
}
