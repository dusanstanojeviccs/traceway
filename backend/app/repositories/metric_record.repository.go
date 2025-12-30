package repositories

import (
	"backend/app/chdb"
	"backend/app/models"
	"context"
	"time"
)

type metricRecordRepository struct{}

func (e *metricRecordRepository) InsertAsync(ctx context.Context, lines []models.MetricRecord) error {
	batch, err := (*chdb.Conn).PrepareBatch(ctx, "INSERT INTO metric_records (name, value, recorded_at)")
	if err != nil {
		return err
	}
	for _, e := range lines {
		if err := batch.Append(e.Name, e.Value, e.RecordedAt); err != nil {
			return err
		}
	}
	return batch.Send()
}

func (e *metricRecordRepository) GetAverageBetween(ctx context.Context, name string, start, end time.Time) (float64, error) {
	var avg float64
	err := (*chdb.Conn).QueryRow(ctx, "SELECT coalesce(avg(value), 0) FROM metric_records WHERE name = ? AND recorded_at >= ? AND recorded_at <= ?", name, start, end).Scan(&avg)
	return avg, err
}

var MetricRecordRepository = metricRecordRepository{}
