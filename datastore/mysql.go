package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const CREATE_USERS_TABLE = `CREATE TABLE IF NOT EXISTS USERS(
	ID int primary key auto_increment,
	USER_NAME text,
	ACCOUNT text,
	AMOUNT int default 0
)`

const CREATE_ORDERS_TABLE = `CREATE TABLE IF NOT EXISTS ORDERS(
	ID int primary key auto_increment,
	ACCOUNT text,
	PRODUCT_NAME text,
	PRICE int,
	ORDER_STATUS text
)`

type sqlDB struct {
	*sql.DB
}

func New() (DB, error) {
	// read the config from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	var (
		username = os.Getenv("SQL_USER")
		password = os.Getenv("SQL_PASSWORD")
		host     = os.Getenv("SQL_HOST")
		dbName   = os.Getenv("SQL_DB")
	)

	// open up our database connection.
	db, err := otelsql.Open("mysql", datasourceName(username, password, host, ""), otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
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

	db, err = otelsql.Open("mysql", datasourceName(username, password, host, dbName))
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
	db.DB.Close()
}

func (db sqlDB) InsertOne(ctx context.Context, p InsertParams) (int64, error) {
	stmt, err := db.PrepareContext(ctx, p.Query)
	if err != nil {
		return 0, fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, p.Vars...)
	if err != nil {
		return 0, fmt.Errorf("statement exec error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("find affected rows error: %w", err)
	}

	return id, nil
}

func (db sqlDB) SelectOne(ctx context.Context, p SelectParams) error {
	stmt, err := db.PrepareContext(ctx, p.Query)
	if err != nil {
		return fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, p.Filters...)
	if err := row.Scan(p.Result...); err != nil {
		return fmt.Errorf("row scan error: %w", err)
	}

	return nil
}

func (db sqlDB) UpdateOne(ctx context.Context, p UpdateParams) error {
	stmt, err := db.PrepareContext(ctx, p.Query)
	if err != nil {
		return fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, p.Vars...); err != nil {
		return fmt.Errorf("statement exec error: %w", err)
	}

	return nil
}

func datasourceName(username, password, host, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbName)
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
