package postgres

import (
	"context"
	"fmt"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	dbUrl, err := url.Parse(config.GetConfig().Database.Address)
	if err != nil {
		panic(err)
	}
	var (
		host    = strings.Split(dbUrl.Host, ":")[0]
		port    = strings.Split(dbUrl.Host, ":")[1]
		user    = dbUrl.User.Username()
		pass, _ = dbUrl.User.Password()
		dbname  = dbUrl.Path[1:]
		mode    = dbUrl.Query().Get("sslmode")
	)
	// compose connection string
	composed := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, mode)
	fmt.Println("composed=", composed)
	db.conn, err = sqlx.Open("postgres", composed)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbConn) Ping(ctx context.Context) error {
	err := db.conn.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
