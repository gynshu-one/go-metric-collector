package postgres

import (
	"context"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/jackc/pgx/v5"
)

type DB interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context) error
}
type face struct {
	conn *pgx.Conn
}

func NewDB() *face {
	return &face{}
}

func (db *face) Connect(ctx context.Context) error {
	var err error
	db.conn, err = pgx.Connect(ctx, config.GetConfig().Database.Address)
	if err != nil {
		return err
	}
	defer db.conn.Close(context.Background())
	return nil
}

func (db *face) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}
