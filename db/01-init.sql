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