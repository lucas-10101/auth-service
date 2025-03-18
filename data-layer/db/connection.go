package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	currentConnection *sql.DB
	driverName        string
	connectionDsn     string
)

type ConnectionWrapper interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func GetConnection() *sql.DB {

	if currentConnection != nil {
		return currentConnection
	}

	conn, err := sql.Open(driverName, connectionDsn) // TODO: move to configuration file
	if err != nil {
		panic("cannot open database connection")
	}

	currentConnection = conn
	return currentConnection
}

func ConfigureConnection(driver, dsn string) {
	driverName = driver
	connectionDsn = dsn
}
