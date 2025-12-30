package repositories

import (
	"backend/app/chdb"
	"backend/app/models"
	"context"
	"time"
)

type transactionRepository struct{}

func (e *transactionRepository) InsertAsync(ctx context.Context, lines []models.Transaction) error {
	batch, err := (*chdb.Conn).PrepareBatch(ctx, "INSERT INTO transactions (id, endpoint, duration, recorded_at, status_code, body_size, client_ip)")
	if err != nil {
		return err
	}
	for _, e := range lines {
		if err := batch.Append(e.Id, e.Endpoint, e.Duration, e.RecordedAt, e.StatusCode, e.BodySize, e.ClientIP); err != nil {
			return err
		}
	}
	return batch.Send()
}

func (e *transactionRepository) CountBetween(ctx context.Context, start, end time.Time) (int64, error) {
	var count uint64
	err := (*chdb.Conn).QueryRow(ctx, "SELECT count() FROM transactions WHERE recorded_at >= ? AND recorded_at <= ?", start, end).Scan(&count)
	return int64(count), err
}

func (e *transactionRepository) FindAll(ctx context.Context, fromDate, toDate time.Time, page, pageSize int, orderBy string) ([]models.Transaction, int64, error) {
	var count uint64
	err := (*chdb.Conn).QueryRow(ctx, "SELECT count() FROM transactions WHERE recorded_at >= ? AND recorded_at <= ?", fromDate, toDate).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	allowedOrderBy := map[string]bool{
		"recorded_at": true,
		"duration":    true,
		"status_code": true,
		"body_size":   true,
	}

	if !allowedOrderBy[orderBy] {
		orderBy = "recorded_at"
	}

	query := "SELECT id, endpoint, duration, recorded_at, status_code, body_size, client_ip FROM transactions WHERE recorded_at >= ? AND recorded_at <= ? ORDER BY " + orderBy + " DESC LIMIT ? OFFSET ?"
	rows, err := (*chdb.Conn).Query(ctx, query, fromDate, toDate, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.Id, &t.Endpoint, &t.Duration, &t.RecordedAt, &t.StatusCode, &t.BodySize, &t.ClientIP); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}

	return transactions, int64(count), nil
}

var TransactionRepository = transactionRepository{}
