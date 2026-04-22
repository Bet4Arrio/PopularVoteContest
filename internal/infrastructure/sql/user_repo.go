package sqlrepo

import (
	"context"
	"errors"

	"github.com/PopularVote/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindUserByPublicID(ctx context.Context, id string) (*user.User, error) {
	const q = `
		SELECT id, public_id, email, "passwordHash", "createAt", "changeAt"
		FROM "user"
		WHERE public_id = $1
	`
	u := &user.User{}
	err := r.db.Pool.QueryRow(ctx, q, id).
		Scan(&u.ID, &u.PublicID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) SaveUser(ctx context.Context, dto *user.CreateUserDTO) (*user.User, error) {
	u := &user.User{
		PublicID: uuid.New().String(),
		Email:    dto.Email,
	}

	const q = `
		INSERT INTO "user" (public_id, email, "passwordHash", "createAt", "changeAt")
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, public_id, email, "passwordHash", "createAt", "changeAt"
	`
	row := r.db.Pool.QueryRow(ctx, q, u.PublicID, dto.Email, dto.Password)
	if err := row.Scan(&u.ID, &u.PublicID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) FindUserByID(ctx context.Context, id string) (*user.User, error) {
	const q = `
		SELECT id, public_id, email, "passwordHash", "createAt", "changeAt"
		FROM "user"
		WHERE public_id = $1
	`
	u := &user.User{}
	err := r.db.Pool.QueryRow(ctx, q, id).
		Scan(&u.ID, &u.PublicID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	const q = `
		SELECT id, public_id, email, "passwordHash", "createAt", "changeAt"
		FROM "user"
		WHERE email = $1
	`
	u := &user.User{}
	err := r.db.Pool.QueryRow(ctx, q, email).
		Scan(&u.ID, &u.PublicID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}
