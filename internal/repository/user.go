package repository

import (
	"context"
	"fmt"
	"time"

	"go-test-api/internal/db"
	"go-test-api/internal/model"

	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository handles user data access
type UserRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

// Upsert creates a new user or updates existing user by email (idempotent operation)
func (r *UserRepository) Upsert(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	now := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}

	user, err := r.queries.UpsertUser(ctx, db.UpsertUserParams{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return &model.UserResponse{
		ID:    fmt.Sprintf("%d", user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// List retrieves all users
func (r *UserRepository) List(ctx context.Context) ([]*model.UserResponse, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	result := make([]*model.UserResponse, len(users))
	for i, user := range users {
		result[i] = &model.UserResponse{
			ID:    fmt.Sprintf("%d", user.ID),
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return result, nil
}
