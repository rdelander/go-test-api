package middleware

import (
	"context"
	"sync"
)

type contextKey string

const dbStatsKey contextKey = "db_stats"

// DBStats tracks database operations for a request
type DBStats struct {
	mu      sync.Mutex
	Queries []QueryInfo
}

// QueryInfo represents a single database query
type QueryInfo struct {
	Query     string
	QueryType string // "SELECT", "INSERT", "UPDATE", "DELETE"
	RowCount  int
}

// NewDBStats creates a new DBStats tracker
func NewDBStats() *DBStats {
	return &DBStats{
		Queries: make([]QueryInfo, 0),
	}
}

// AddQuery records a query execution
func (s *DBStats) AddQuery(queryType, query string, rowCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Queries = append(s.Queries, QueryInfo{
		Query:     query,
		QueryType: queryType,
		RowCount:  rowCount,
	})
}

// Summary returns aggregated statistics
func (s *DBStats) Summary() (total, selects, inserts, updates, deletes int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	total = len(s.Queries)
	for _, q := range s.Queries {
		switch q.QueryType {
		case "SELECT":
			selects++
		case "INSERT":
			inserts++
		case "UPDATE":
			updates++
		case "DELETE":
			deletes++
		}
	}
	return
}

// GetDBStats retrieves DBStats from context
func GetDBStats(ctx context.Context) *DBStats {
	if stats, ok := ctx.Value(dbStatsKey).(*DBStats); ok {
		return stats
	}
	return nil
}

// WithDBStats adds DBStats to context
func WithDBStats(ctx context.Context) context.Context {
	return context.WithValue(ctx, dbStatsKey, NewDBStats())
}
