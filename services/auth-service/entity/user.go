package entity

import (
	"context"
	"time"
)

const (
	RoleStudent  = "student"
	RoleLecturer = "lecturer"
)

type User struct {
	ID             uint       `gorm:"primary key" json:"id"`
	Username       string     `gorm:"unique;not null" json:"username"`
	Email          string     `gorm:"unique;not null" json:"email"`
	Password       string     `gorm:"not null"`
	FullName       string     `gorm:"type:varchar(150)" json:"full_name"`
	Role           string     `gorm:"type:varchar(20); not null; default:'student'" json:"-"`
	Bio            string     `gorm:"type:text" json:"bio"`
	ProfilePicture string     `gorm:"type:text" json:"profile_picture"`
	GoogleID       *string    `gorm:"unique;type:varchar(255)" json:"google_id"`
	IsVerified     bool       `gorm:"default:false" json:"is_verified"`
	DOB            *time.Time `gorm:"type:date" json:"dob"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByGoogle(ctx context.Context, googleID string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type CacheRepository interface {
	SetOTP(ctx context.Context, email, code string) error
	GetOTP(ctx context.Context, email string) (string, error)
	DeleteOTP(ctx context.Context, email string) error
	SetRefreshToken(ctx context.Context, userID uint, token string) error
	GetRefreshToken(ctx context.Context, userID uint) (string, error)
	DeleteRefreshToken(ctx context.Context, userID uint) error
}
