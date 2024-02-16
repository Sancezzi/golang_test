package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

type DbManager struct {
	db  *reform.DB
	sql *sql.DB
}

func Connect(config DbConfig) (*DbManager, error) {
	log.Println("INFO: Connecting to DB...")

	sql, err := sql.Open("pgx", getConnectStr(config))
	if err != nil {
		return nil, err
	}

	if err = sql.Ping(); err != nil {
		return nil, err
	}

	db := reform.NewDB(sql, postgresql.Dialect, nil)

	return &DbManager{
		db:  db,
		sql: sql,
	}, nil
}

func (mgr *DbManager) GetConnection() *reform.DB {
	return mgr.db
}

func (mgr *DbManager) Close() error {
	if err := mgr.sql.Close(); err != nil {
		log.Println("INFO: Cannot close DB connection")
		return err
	}
	return nil
}

func getConnectStr(config DbConfig) string {
	// postgres://username:password@localhost:5432/database_name

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Username,
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.DbName,
	)
}
