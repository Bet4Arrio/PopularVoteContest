package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, name, email, plainPassword string) (*User, error) {
	existing, _ := s.repo.FindUserByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailTaken
	}

	if email == "" {
		return nil, ErrInvalidEmail
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if plainPassword == "" {
		return nil, ErrInvalidPassword
	}
	// TODO>  use a better password hashing strategy with a unique salt per user, e.g. argon2id.e
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	dto := &CreateUserDTO{
		Name:     name,
		Email:    email,
		Password: string(passwordHash),
	}
	return s.repo.SaveUser(ctx, dto)
}

func (s *Service) GetUserByPublicID(ctx context.Context, id string) (*User, error) {
	return s.repo.FindUserByPublicID(ctx, id)
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.repo.FindUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindUserByEmail(ctx, email)
}

func (s *Service) AuthenticateUser(ctx context.Context, email, plainPassword string) (*User, error) {
	u, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword)); err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}
