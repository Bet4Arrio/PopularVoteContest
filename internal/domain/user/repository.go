package user

import "context"

type Repository interface {
	SaveUser(ctx context.Context, dto *CreateUserDTO) (User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}
