package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailTaken      = errors.New("email already in use")
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidName     = errors.New("invalid name")
)

// User is the aggregate root for the user domain.
type User struct {
	ID           int       `bson:"_id,omitempty" json:"id"`
	PublicID     string    `bson:"public_id"     json:"public_id"` // UUID for external reference
	Email        string    `bson:"email"         json:"email"`
	PasswordHash string    `bson:"password_hash" json:"-"`
	Role         Role      `bson:"role"          json:"role"`
	CreatedAt    time.Time `bson:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"    json:"updated_at"`
}

func NewUser(email, passwordHash string, role Role) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if len(passwordHash) < 8 {
		return nil, ErrInvalidPassword
	}
	if role != RoleUser && role != RoleAdmin {
		return nil, ErrInvalidRole
	}

	return &User{
		PublicID:     uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// UpdateEmail updates the user's email.
