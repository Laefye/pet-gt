package repository

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	Sessions []Session `gorm:"foreignKey:UserID"`
}

type CreateUserRequest struct {
	Username string
	Password string
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	user := &User{
		ID:       ulid.Make().String(),
		Username: req.Username,
		Password: req.Password,
	}
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where(&User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
