package adapters

import (
	"context"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"time"
)

type DBAdapter interface {
	StoreMetrics(context.Context, []*entity.Metrics) error
	GetMetrics(context.Context) ([]*entity.Metrics, error)
}
type dbAdapter struct {
	conn *sqlx.DB
	DBAdapter
}

// NewAdapter creates new dbAdapter (Context is for committing initial scheme)
func NewAdapter(ctx context.Context, conn *sqlx.DB) *dbAdapter {
	adap := &dbAdapter{conn: conn}
	c, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	err := adap.commitScheme(c)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to commit initial scheme")
	}
	return adap
}

const schema = `
CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    delta BIGINT,
    value FLOAT,
    hash VARCHAR(255)
);
`
const insert = `
INSERT INTO metrics ( id, type, delta, value, hash) 
VALUES ($1, $2, $3, $4, $5)
`
const selectAll = `
SELECT * FROM metrics
`
const insertOrUpdate = `
INSERT INTO metrics ( id, type, delta, value, hash)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET type = $2, delta = $3, value = $4, hash = $5
`

func (a *dbAdapter) StoreMetrics(ctx context.Context, metrics []*entity.Metrics) error {
	// ! I don't know how to do it better. Define timout on low level functions or on high level?
	c, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	tx, err := a.conn.BeginTx(c, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to begin transaction StoreMetrics")
		return err
	}
	defer tx.Rollback()
	for _, m := range metrics {
		_, err = tx.ExecContext(c, insertOrUpdate, m.ID, m.MType, m.Delta, m.Value, m.Hash)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Unable to commit transaction StoreMetrics")
		return err
	}
	return nil
}
func (a *dbAdapter) GetMetrics(ctx context.Context) ([]*entity.Metrics, error) {
	c, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	rows, err := a.conn.QueryxContext(c, selectAll)
	if err != nil || rows.Err() != nil {
		log.Error().Err(err).Msg("Unable to get metrics")
		return nil, err
	}
	defer rows.Close()
	metrics := make([]*entity.Metrics, 0)
	for rows.Next() {
		var m entity.Metrics
		err = rows.StructScan(&m)
		if err != nil {
			log.Error().Err(err).Msg("Unable to scan and unmarshal metrics from db")
			return nil, err
		}
		metrics = append(metrics, &m)
	}
	return metrics, nil
}
func (a *dbAdapter) commitScheme(ctx context.Context) error {
	_, err := a.conn.ExecContext(ctx, schema)
	if err != nil {
		return err
	}
	tx, err := a.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, insert, "test2", "test1", 1, 1.0, "test1")
	if err != nil {
		return err
	}
	err = tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}
