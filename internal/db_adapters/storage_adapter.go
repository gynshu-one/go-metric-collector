package db_adapters

import (
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBAdapter interface {
	StoreMetrics([]*entity.Metrics) error
	GetMetrics() ([]*entity.Metrics, error)
}

type dbAdapter struct {
	conn *sqlx.DB
	DBAdapter
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

func (a *dbAdapter) StoreMetrics(metrics []*entity.Metrics) error {
	tx := a.conn.MustBegin()
	defer tx.Rollback()
	for _, m := range metrics {
		mt, err := a.getMetricsByID(m.ID)
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
func (a *dbAdapter) GetMetrics() ([]*entity.Metrics, error) {
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
func (a *dbAdapter) getMetricsByID(id string) (*entity.Metrics, error) {
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
