package repository

import (
	"context"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Session struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	UserAgent string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`

	User *User `gorm:"foreignKey:UserID"`
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, userID, userAgent string) (*Session, error) {
	session := &Session{
		ID:        ulid.Make().String(),
		UserID:    userID,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*Session, error) {
	var session Session
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}
