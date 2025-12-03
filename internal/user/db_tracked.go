package user

import (
	"context"
	"strings"

	"go-test-api/internal/middleware"
	"go-test-api/internal/user/db"
)

// TrackedQueries wraps db.Queries to track operations
type TrackedQueries struct {
	*db.Queries
}

// NewTrackedQueries creates a queries wrapper that tracks operations
func NewTrackedQueries(q *db.Queries) *TrackedQueries {
	return &TrackedQueries{Queries: q}
}

func trackQuery(ctx context.Context, queryType, query string, rowCount int) {
	if stats := middleware.GetDBStats(ctx); stats != nil {
		stats.AddQuery(queryType, query, rowCount)
	}
}

func inferQueryType(query string) string {
	query = strings.ToUpper(strings.TrimSpace(query))
	if strings.HasPrefix(query, "SELECT") {
		return "SELECT"
	} else if strings.HasPrefix(query, "INSERT") {
		return "INSERT"
	} else if strings.HasPrefix(query, "UPDATE") {
		return "UPDATE"
	} else if strings.HasPrefix(query, "DELETE") {
		return "DELETE"
	}
	return "OTHER"
}

func (q *TrackedQueries) ListUsers(ctx context.Context) ([]db.User, error) {
	users, err := q.Queries.ListUsers(ctx)
	if err == nil {
		trackQuery(ctx, "SELECT", "ListUsers", len(users))
	}
	return users, err
}

func (q *TrackedQueries) ListUsersByEmail(ctx context.Context, email string) ([]db.User, error) {
	users, err := q.Queries.ListUsersByEmail(ctx, email)
	if err == nil {
		trackQuery(ctx, "SELECT", "ListUsersByEmail: "+email, len(users))
	}
	return users, err
}

func (q *TrackedQueries) UpsertUser(ctx context.Context, arg db.UpsertUserParams) (db.User, error) {
	user, err := q.Queries.UpsertUser(ctx, arg)
	if err == nil {
		trackQuery(ctx, "INSERT", "UpsertUser: "+arg.Email, 1)
	}
	return user, err
}
