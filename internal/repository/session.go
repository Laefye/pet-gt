package repository

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Session struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	UserAgent string    `gorm:"not null"`
	Time      time.Time `gorm:"not null"`

	User *User `gorm:"foreignKey:UserID"`
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

type CreateSessionRequest struct {
	UserID    string
	UserAgent string
}

func (r *SessionRepository) CreateSession(ctx context.Context, req CreateSessionRequest) (*Session, error) {
	session := &Session{
		ID:        ulid.Make().String(),
		UserID:    req.UserID,
		UserAgent: req.UserAgent,
		Time:      time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	var session Session
	if err := r.db.WithContext(ctx).Preload("User").Where(&Session{ID: id}).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}
