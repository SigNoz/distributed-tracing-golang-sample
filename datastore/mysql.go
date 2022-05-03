package datastore

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "password"
	hostname = "localhost:3306"
	dbName   = "signoz"
)

type sqlDB struct {
	*sql.DB
}

func New() (DB, error) {
	// open up our database connection.
	db, err := sql.Open("mysql", datasourceName(""))
	if err != nil {
		return nil, fmt.Errorf("open main db error: %w", err)
	}
	defer db.Close()

	// create signoz db
	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName); err != nil {
		return nil, fmt.Errorf("signoz db create error: %w", err)
	}

	// close the exising connection. db.Close() is idempotent. Hence, it is safe to close the db here.
	db.Close()

	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		return nil, fmt.Errorf("open signoz db error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	log.Printf("Successfully connected to %s DB\n", dbName)

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("create tables error: %w", err)
	}

	return sqlDB{db}, nil
}

func (db sqlDB) Close() {
	db.Close()
}

func (db sqlDB) InsertOne(p InsertParams) (int64, error) {
	stmt, err := db.Prepare(p.Query)
	if err != nil {
		return 0, fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.Vars...)
	if err != nil {
		return 0, fmt.Errorf("statement exec error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("find affected rows error: %w", err)
	}

	return id, nil
}

func (db sqlDB) SelectOne(p SelectParams) error {
	stmt, err := db.Prepare(p.Query)
	if err != nil {
		return fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(p.Filters...)
	if err := row.Scan(p.Result...); err != nil {
		return fmt.Errorf("row scan error: %w", err)
	}

	return nil
}

func (db sqlDB) UpdateOne(p UpdateParams) error {
	stmt, err := db.Prepare(p.Query)
	if err != nil {
		return fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(p.Vars...); err != nil {
		return fmt.Errorf("statement exec error: %w", err)
	}

	return nil
}

func datasourceName(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func createTables(db *sql.DB) error {
	if _, err := db.Exec(CREATE_USERS_TABLE); err != nil {
		return fmt.Errorf("create user table error: %w", err)
	}

	if _, err := db.Exec(CREATE_ORDERS_TABLE); err != nil {
		return fmt.Errorf("create orders table error: %w", err)
	}

	return nil
}
