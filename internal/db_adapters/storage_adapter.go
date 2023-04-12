package db_adapters

import (
	"context"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type DbAdapter interface {
	StoreMetrics(context.Context, []*entity.Metrics) error
	GetMetrics(context.Context) ([]*entity.Metrics, error)
	Test() bool
}

type dbAdapter struct {
	conn *sqlx.DB
	DbAdapter
}

func NewAdapter(conn *sqlx.DB) *dbAdapter {
	adap := &dbAdapter{conn: conn}
	err := adap.commitScheme()
	if err != nil {
		log.Fatal(err)
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
const update = `
UPDATE metrics SET type = $2, delta = $3, value = $4, hash = $5 WHERE id = $1
`
const getByID = `
SELECT * FROM metrics WHERE id = $1
`

const test = `
SELECT true
FROM  metrics  AS tbl
WHERE
    tbl::text LIKE $1
LIMIT 1`

func (a *dbAdapter) Test() bool {
	rows, err := a.conn.Queryx(test, "%PopulateCounter1625326%")
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}
func (a *dbAdapter) StoreMetrics(ctx context.Context, metrics []*entity.Metrics) error {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	tx := a.conn.MustBegin()
	defer tx.Rollback()
	for _, m := range metrics {
		mt, err := a.getMetricsByName(ctx, m.ID)
		if err != nil {
			return err
		}
		if mt.ID != "" {
			tx.MustExec(update, mt.ID, m.MType, m.Delta, m.Value, m.Hash)
			continue
		}
		tx.MustExec(insert, m.ID, m.MType, m.Delta, m.Value, m.Hash)
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (a *dbAdapter) GetMetrics(ctx context.Context) ([]*entity.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	rows, err := a.conn.Queryx(selectAll)
	defer rows.Close()
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
func (a *dbAdapter) commitScheme() error {
	a.conn.MustExec(schema)
	tx := a.conn.MustBegin()
	defer tx.Rollback()
	tx.MustExec(insert, "test2", "test1", 1, 1.0, "test1")
	err := tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}
func (a *dbAdapter) getMetricsById(ctx context.Context, id int64) (*entity.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	rows, err := a.conn.Queryx(getByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var m entity.Metrics
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
	}
	return &m, nil
}
func (a *dbAdapter) getMetricsByName(ctx context.Context, id string) (*entity.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	rows, err := a.conn.Queryx(getByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var m entity.Metrics
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
	}
	return &m, nil
}
func (a *dbAdapter) updateMetrics(ctx context.Context, metrics []*entity.Metrics) error {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	tx := a.conn.MustBegin()
	for _, m := range metrics {
		tx.MustExec(update, m.MType, m.Delta, m.Value, m.Hash, m.ID)
	}
	defer tx.Rollback()
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}