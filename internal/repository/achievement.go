package repository

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type AchievementName string

const (
	AchievementFirstLogin = AchievementName("first_login")
)

func (a AchievementName) IsValid() bool {
	switch a {
	case AchievementFirstLogin:
		return true
	default:
		return false
	}
}

func (a AchievementName) String() string {
	switch a {
	case AchievementFirstLogin:
		return "First Login"
	default:
		return "Unknown Achievement"
	}
}

func (a AchievementName) ImageURL() string {
	switch a {
	case AchievementFirstLogin:
		return "/public/img/first.png"
	default:
		return ""
	}
}

type Achievement struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"index,not null"`
	Name      string    `gorm:"not null"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"not null"`
}

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

type CreateAchievementRequest struct {
	UserID string
	Name   AchievementName
}

func (r *AchievementRepository) Create(ctx context.Context, req *CreateAchievementRequest) (*Achievement, error) {
	achievement := &Achievement{
		ID:        ulid.Make().String(),
		UserID:    req.UserID,
		Name:      string(req.Name),
		CreatedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(achievement).Error; err != nil {
		return nil, err
	}
	return achievement, nil
}

func (r *AchievementRepository) GetByUserID(ctx context.Context, userID string) ([]*Achievement, error) {
	var achievements []*Achievement
	err := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *AchievementRepository) Contains(ctx context.Context, userID string, name AchievementName) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Achievement{}).Where("user_id = ? AND name = ?", userID, string(name)).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
