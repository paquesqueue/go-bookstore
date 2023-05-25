package api

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/paquesqueue/bookstore/common"
)

const sqlQuries = `
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY NOT NULL, 
    title TEXT NOT NULL, 
    authors TEXT[] NOT NULL, 
    publisher TEXT NOT NULL, 
    isbn TEXT NOT NULL, 
    price BIGINT NOT NULL, 
    quantity BIGINT NOT NULL, 
    created_by TEXT NOT NULL, 
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
    );

CREATE TABLE IF NOT EXISTS users (
	username TEXT PRIMARY KEY NOT NULL UNIQUE, 
	email TEXT NOT NULL UNIQUE,
	fullname TEXT NOT NULL,
	hashed_password TEXT NOT NULL,	
	created_at TIMESTAMP DEFAULT NOW() NOT NULL
);
`

func InitDB(c common.Config, log common.Log) (*sql.DB, error) {
	db, err := sql.Open(c.DriverName, c.Url)
	if err != nil {
		log.Errorf("Error Open Database : %v", err)
	}

	_, err = db.Exec(sqlQuries)
	if err != nil {
		log.Errorf("Error Create Tables : %v", err)
	}
	return db, err
}

type Query struct {
	*sql.DB
}

func NewDB(db *sql.DB) Query {
	return Query{db}
}
