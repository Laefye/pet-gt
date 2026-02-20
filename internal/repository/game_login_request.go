package repository

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GameLoginRequest struct {
	ID               string  `gorm:"primaryKey"`
	AuthorizedUserID *string `gorm:"index;column:authorized_user_id"`
	AuthorizedUser   *User   `gorm:"foreignKey:AuthorizedUserID;references:ID"`
}

type GameLoginRequestRepository struct {
	db *gorm.DB
}

func NewGameLoginRequestRepository(db *gorm.DB) *GameLoginRequestRepository {
	return &GameLoginRequestRepository{db: db}
}

func (r *GameLoginRequestRepository) Create(ctx context.Context) (*GameLoginRequest, error) {
	req := &GameLoginRequest{
		ID: ulid.Make().String(),
	}
	if err := r.db.WithContext(ctx).Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

func (r *GameLoginRequestRepository) GetByID(ctx context.Context, id string) (*GameLoginRequest, error) {
	var req GameLoginRequest
	err := r.db.WithContext(ctx).Preload("AuthorizedUser").Where("id = ?", id).First(&req).Error
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
