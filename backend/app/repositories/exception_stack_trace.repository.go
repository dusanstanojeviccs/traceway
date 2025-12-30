package repositories

import (
	"backend/app/chdb"
	"backend/app/models"
	"context"
	"time"
)

type exceptionStackTraceRepository struct{}

func (e *exceptionStackTraceRepository) InsertAsync(ctx context.Context, lines []models.ExceptionStackTrace) error {
	batch, err := (*chdb.Conn).PrepareBatch(ctx, "INSERT INTO exception_stack_traces (transaction_id, exception_hash, stack_trace, recorded_at)")
	if err != nil {
		return err
	}
	for _, e := range lines {
		if err := batch.Append(e.TransactionId, e.ExceptionHash, e.StackTrace, e.RecordedAt); err != nil {
			return err
		}
	}
	return batch.Send()
}

func (e *exceptionStackTraceRepository) CountBetween(ctx context.Context, start, end time.Time) (int64, error) {
	var count uint64
	err := (*chdb.Conn).QueryRow(ctx, "SELECT count() FROM exception_stack_traces WHERE recorded_at >= ? AND recorded_at <= ?", start, end).Scan(&count)
	return int64(count), err
}

func (e *exceptionStackTraceRepository) FindGrouped(ctx context.Context, fromDate, toDate time.Time, page, pageSize int, orderBy string) ([]models.ExceptionGroup, int64, error) {
	// Count unique hashes
	var count uint64
	err := (*chdb.Conn).QueryRow(ctx, "SELECT uniq(exception_hash) FROM exception_stack_traces WHERE recorded_at >= ? AND recorded_at <= ?", fromDate, toDate).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	allowedOrderBy := map[string]bool{
		"last_seen":  true,
		"first_seen": true,
		"count":      true,
	}

	if !allowedOrderBy[orderBy] {
		orderBy = "count"
	}

	// Re-writing query string construction for OrderBy to avoid binding issue
	fullQuery := "SELECT exception_hash, any(stack_trace), max(recorded_at) as last_seen, min(recorded_at) as first_seen, count() as count FROM exception_stack_traces WHERE recorded_at >= ? AND recorded_at <= ? GROUP BY exception_hash ORDER BY " + orderBy + " DESC LIMIT ? OFFSET ?"

	rows, err := (*chdb.Conn).Query(ctx, fullQuery, fromDate, toDate, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []models.ExceptionGroup
	for rows.Next() {
		var g models.ExceptionGroup
		if err := rows.Scan(&g.ExceptionHash, &g.StackTrace, &g.LastSeen, &g.FirstSeen, &g.Count); err != nil {
			return nil, 0, err
		}
		groups = append(groups, g)
	}

	return groups, int64(count), nil
}

var ExceptionStackTraceRepository = exceptionStackTraceRepository{}
