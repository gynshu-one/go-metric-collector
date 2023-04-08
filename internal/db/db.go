package db

import (
	"context"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/jackc/pgx/v5"
)

var Db Face

type Face struct {
	conn *pgx.Conn
	Get  Get
	Do   Do
}
type Do interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context) error
}
type Get interface {
	Conn() *pgx.Conn
}

func NewDb() *Face {
	return &Face{}
}

func (db *Face) Connect(ctx context.Context) error {
	var err error
	db.conn, err = pgx.Connect(ctx, configs.CFG.DbAddress)
	if err != nil {
		return err
	}
	defer db.conn.Close(context.Background())
	return nil
}

func (db *Face) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}
func (db *Face) Conn() *pgx.Conn {
	return db.conn
}
