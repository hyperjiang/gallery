package provider

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hyperjiang/gallery-service/app/config"

	"github.com/jmoiron/sqlx"
	_ "gopkg.in/go-sql-driver/mysql.v1" // mysql driver
)

// try to open DB connection, this function is intended for test mocking
var tryOpenDB func(URI string) (*sqlx.DB, error)

// openDB - open a DB connection
func openDB(config *config.DatabaseSettings) (*sqlx.DB, error) {

	if tryOpenDB == nil { // not defined, use default function
		tryOpenDB = func(URI string) (*sqlx.DB, error) {
			return sqlx.Connect("mysql", URI)
		}
	}

	db, err := tryOpenDB(config.URI)

	if err != nil {
		log.Println("Cannot open DB")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("Cannot ping DB")
		db.Close()
		return nil, err
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	return db, nil
}

// GenerateInsertSQL - generate SQL for insert
func GenerateInsertSQL(table string, data interface{}) (string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return "", errors.New("data must be a pointer of table struct")
	}
	v = v.Elem()

	var fields, placeholders []string
	for i := 0; i < v.NumField(); i++ {
		column := v.Type().Field(i).Tag.Get("db")
		fields = append(fields, column)
		placeholders = append(placeholders, ":"+column)
	}

	sql := "INSERT INTO " + table + "(" + strings.Join(fields, ",") +
		") VALUES (" + strings.Join(placeholders, ",") + ")"

	return sql, nil
}

// DBInsert - helper for inserting data into db table
func DBInsert(db *sqlx.DB, table string, data interface{}) (sql.Result, error) {
	sql, err := GenerateInsertSQL(table, data)
	if err != nil {
		return nil, err
	}

	return db.NamedExec(sql, data)
}

// TxInsert - helper for inserting data into db table (using transaction)
func TxInsert(tx *sqlx.Tx, table string, data interface{}) (sql.Result, error) {
	sql, err := GenerateInsertSQL(table, data)
	if err != nil {
		return nil, err
	}

	return tx.NamedExec(sql, data)
}
