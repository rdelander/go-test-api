package database

import (
	"context"
	"strings"

	"go-test-api/internal/middleware"

	"github.com/jackc/pgx/v5"
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
		queryType := inferQueryType(sql)
		rowCount := int(data.CommandTag.RowsAffected())
		stats.AddQuery(queryType, sql, rowCount)
	}
}

func inferQueryType(sql string) string {
	sql = strings.TrimSpace(strings.ToUpper(sql))
	if strings.HasPrefix(sql, "SELECT") {
		return "SELECT"
	} else if strings.HasPrefix(sql, "INSERT") {
		return "INSERT"
	} else if strings.HasPrefix(sql, "UPDATE") {
		return "UPDATE"
	} else if strings.HasPrefix(sql, "DELETE") {
		return "DELETE"
	}
	return "OTHER"
}
