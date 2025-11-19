package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-test-api/internal/model"
)

// UserRepository handles user data access
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	query := `
		INSERT INTO users (name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, created_at
	`

	var user model.UserResponse
	var createdAt time.Time

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Name,
		req.Email,
		time.Now(),
		time.Now(),
	).Scan(&user.ID, &user.Name, &user.Email, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// Upsert creates a new user or updates existing user by email (idempotent operation)
func (r *UserRepository) Upsert(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	query := `
		INSERT INTO users (name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) 
		DO UPDATE SET 
			name = EXCLUDED.name,
			updated_at = EXCLUDED.updated_at
		RETURNING id, name, email
	`

	var user model.UserResponse

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Name,
		req.Email,
		time.Now(),
		time.Now(),
	).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.UserResponse, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	var user model.UserResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// List retrieves all users
func (r *UserRepository) List(ctx context.Context) ([]*model.UserResponse, error) {
	query := `SELECT id, name, email FROM users ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*model.UserResponse
	for rows.Next() {
		var user model.UserResponse
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
