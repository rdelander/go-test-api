package database

import (
	"context"

	"go-test-api/internal/middleware"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type tracerContextKey string

const sqlKey tracerContextKey = "sql"

// QueryTracer implements pgx.QueryTracer to track database operations
type QueryTracer struct{}

func (t *QueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// Store SQL in context for later retrieval
	return context.WithValue(ctx, sqlKey, data.SQL)
}

func (t *QueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if stats := middleware.GetDBStats(ctx); stats != nil {
		sql, _ := ctx.Value(sqlKey).(string)
		queryType := inferQueryType(data.CommandTag)
		rowCount := int(data.CommandTag.RowsAffected())
		stats.AddQuery(queryType, sql, rowCount)
	}
}

func inferQueryType(tag pgconn.CommandTag) string {
	switch {
	case tag.Select():
		return "SELECT"
	case tag.Insert():
		return "INSERT"
	case tag.Update():
		return "UPDATE"
	case tag.Delete():
		return "DELETE"
	default:
		return "OTHER"
	}
}
