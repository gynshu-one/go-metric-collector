package db_adapters

import (
	"context"
	"github.com/google/uuid"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/jmoiron/sqlx"
	"time"
)

type DbAdapter interface {
	StoreMetrics(context.Context, []*entity.Metrics) error
	GetMetrics(context.Context) ([]*entity.Metrics, error)
	CommitScheme() error
}

type dbAdapter struct {
	conn *sqlx.DB
	DbAdapter
}

func NewAdapter(conn *sqlx.DB) *dbAdapter {
	return &dbAdapter{conn: conn}
}

const schema = `
CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    delta BIGINT,
    value FLOAT,
    hash VARCHAR(255)
);
`
const insert = `
INSERT INTO metrics (id, name, type, delta, value, hash) VALUES ($1, $2, $3, $4, $5, $6)
`
const selectAll = `
SELECT * FROM metrics
`

func (a dbAdapter) StoreMetrics(ctx context.Context, metrics []*entity.Metrics) error {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	tx := a.conn.MustBegin()
	for _, m := range metrics {
		tx.MustExec(insert, uuid.New().String(), m.MType, m.Delta, m.Value, m.Hash)
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (a dbAdapter) GetMetrics(ctx context.Context) ([]*entity.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	rows, err := a.conn.Queryx(selectAll)
	if err != nil {
		return nil, err
	}
	metrics := make([]*entity.Metrics, 0)
	for rows.Next() {
		var m entity.Metrics
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &m)
	}
	return metrics, nil
}
func (a dbAdapter) CommitScheme() error {
	a.conn.MustExec(schema)
	tx := a.conn.MustBegin()
	tx.MustExec(insert, "test", "test", 1, 1.0, "test")
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
