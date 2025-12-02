package user

import (
	"context"
	"fmt"
	"time"

	"go-test-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// Repo defines the interface for user data access
type Repo interface {
	Upsert(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	List(ctx context.Context) ([]*UserResponse, error)
	ListByEmail(ctx context.Context, email string) ([]*UserResponse, error)
}

// Repository handles user data access
type Repository struct {
	queries *db.Queries
}

// NewRepository creates a new Repository
func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

// Upsert creates or updates a user by email
func (r *Repository) Upsert(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	user, err := r.queries.UpsertUser(ctx, db.UpsertUserParams{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}
	return &UserResponse{ID: fmt.Sprintf("%d", user.ID), Name: user.Name, Email: user.Email}, nil
}

// List returns all users
func (r *Repository) List(ctx context.Context) ([]*UserResponse, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	res := make([]*UserResponse, len(users))
	for i, u := range users {
		res[i] = &UserResponse{ID: fmt.Sprintf("%d", u.ID), Name: u.Name, Email: u.Email}
	}
	return res, nil
}

// ListByEmail filters users by email (repository adds SQL wildcard)
func (r *Repository) ListByEmail(ctx context.Context, email string) ([]*UserResponse, error) {
	pattern := "%" + email + "%"
	users, err := r.queries.ListUsersByEmail(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list users by email: %w", err)
	}
	res := make([]*UserResponse, len(users))
	for i, u := range users {
		res[i] = &UserResponse{ID: fmt.Sprintf("%d", u.ID), Name: u.Name, Email: u.Email}
	}
	return res, nil
}
