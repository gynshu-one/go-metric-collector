package postgres

import (
	"context"
	"fmt"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

type DBConn interface {
	GetConn() *sqlx.DB
	Connect() error
	Ping(ctx context.Context) error
}
type dbConn struct {
	conn *sqlx.DB
}

func NewDB() *dbConn {
	return &dbConn{}
}

func (db *dbConn) GetConn() *sqlx.DB {
	return db.conn
}
func (db *dbConn) Connect() error {
	dbURL, err := url.Parse(config.GetConfig().Database.Address)
	if err != nil {
		panic(err)
	}
	var (
		host    = strings.Split(dbURL.Host, ":")[0]
		port    = strings.Split(dbURL.Host, ":")[1]
		user    = dbURL.User.Username()
		pass, _ = dbURL.User.Password()
		dbname  = dbURL.Path[1:]
		mode    = dbURL.Query().Get("sslmode")
	)
	composed := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, mode)
	log.Debug().Msgf("Connecting to database: %s", composed)
	db.conn, err = sqlx.Open("postgres", composed)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbConn) Ping(ctx context.Context) error {
	err := db.conn.PingContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Ping to database failed")
		return err
	}
	return nil
}
