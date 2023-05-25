package api

import (
	"github.com/lib/pq"
)

func (db Query) InsertBook(req RequestBook) (*ResponseBook, error) {
	const query = `INSERT INTO books 
	(title, authors, publisher, isbn, price, quantity, created_by) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(req.Title, pq.Array(req.Authors), req.Publisher, req.Isbn, req.Price, req.Quantity, req.Created_by)

	resp := &ResponseBook{}
	err = row.Scan(&resp.Id, &resp.Title, pq.Array(&resp.Authors), &resp.Publisher, &resp.Isbn, &resp.Price, &resp.Quantity, &resp.Created_by, &resp.Created_at)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (db Query) SelectAllBooks(params GetAllParams) ([]ResponseBook, error) {
	const query = `SELECT id, title, authors, publisher, isbn, price, quantity, created_by, created_at 
	FROM books
	ORDER BY id
	LIMIT $1
	OFFSET $2;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := []ResponseBook{}
	for rows.Next() {
		result := &ResponseBook{}
		err = rows.Scan(&result.Id, &result.Title, pq.Array(&result.Authors), &result.Publisher, &result.Isbn, &result.Price, &result.Quantity, &result.Created_by, &result.Created_at)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *result)
	}
	return resp, nil

}

func (db Query) SelectBookByID(id uint64) (*ResponseBook, error) {
	const query = `SELECT * FROM books WHERE id = $1;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	resp := &ResponseBook{}
	err = row.Scan(&resp.Id, &resp.Title, pq.Array(&resp.Authors), &resp.Publisher, &resp.Isbn, &resp.Price, &resp.Quantity, &resp.Created_by, &resp.Created_at)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (db Query) UpdateBook(id uint64, req RequestBook) (*ResponseBook, error) {
	const query = `UPDATE books 
	SET title = $1, authors = $2, publisher = $3, isbn = $4, price = $5, quantity = $6, created_by = $7 
	WHERE id = $8 
	RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(req.Title, pq.Array(req.Authors), req.Publisher, req.Isbn, req.Price, req.Quantity, req.Created_by, id)
	resp := &ResponseBook{}
	err = row.Scan(&resp.Id, &resp.Title, pq.Array(&resp.Authors), &resp.Publisher, &resp.Isbn, &resp.Price, &resp.Quantity, &resp.Created_by, &resp.Created_at)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (db Query) DeleteBook(id uint64) error {
	const query = `DELETE FROM books WHERE id = $1;`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
