//go:build unit

package api

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsertBook(t *testing.T) {
	t.Run("TestInsertShouldReturnNoError", func(t *testing.T) {
		// Arrange
		mockData := RequestBook{
			Title:      "mockTitle",
			Authors:    []string{"mockAuthor A", "mockAuthor B"},
			Publisher:  "mockPublisher",
			Isbn:       "1234567890",
			Price:      1000,
			Quantity:   100,
			Created_by: "mockAdmin",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mockCreated_at := time.Now()
		row := sqlmock.NewRows([]string{"id", "title", "authors", "publisher", "isbn", "price", "quantity", "created_by", "created_at"}).
			AddRow(1, mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO books (title, authors, publisher, isbn, price, quantity, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`))
		get.ExpectQuery().
			WithArgs(mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by).
			WillReturnRows(row)

		query := NewDB(db)

		// Act
		result, err := query.InsertBook(mockData)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, mockData.Title, result.Title)
		assert.Equal(t, mockData.Authors, result.Authors)
		assert.Equal(t, mockData.Publisher, result.Publisher)
		assert.Equal(t, mockData.Isbn, result.Isbn)
		assert.Equal(t, mockData.Price, result.Price)
		assert.Equal(t, mockData.Quantity, result.Quantity)
		assert.Equal(t, uint64(1), result.Id)
		assert.Equal(t, mockCreated_at, result.Created_at)
	})

	t.Run("TestInsertShouldReturnError", func(t *testing.T) {
		// Arrange
		mockData := RequestBook{}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		get := mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO books (title, authors, publisher, isbn, price, quantity, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`))
		get.ExpectQuery().
			WithArgs(mockData).
			WillReturnError(&pq.Error{Message: "db connection error"})

		query := NewDB(db)

		// Act
		result, err := query.InsertBook(mockData)

		// Assert
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestSelectAllBooks(t *testing.T) {
	t.Run("TestSelectAllBooksShouldReturnNoError", func(t *testing.T) {
		// Arrange
		mockData := []RequestBook{
			{
				Title:      "mockTitle 1",
				Authors:    []string{"mockAuthor A", "mockAUthor B"},
				Publisher:  "mockPublsiher",
				Isbn:       "1234567890",
				Price:      1000,
				Quantity:   100,
				Created_by: "mockAdmin",
			},
			{
				Title:      "mockTitle 2",
				Authors:    []string{"mockAuthor A", "mockAUthor B"},
				Publisher:  "mockPublsiher",
				Isbn:       "1234567890",
				Price:      1000,
				Quantity:   100,
				Created_by: "mockAdmin",
			},
		}
		params := GetAllParams{
			Limit:  1,
			Offset: 2,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mockCreated_at := time.Now()
		row := sqlmock.NewRows([]string{"id", "title", "authors", "publisher", "isbn", "price", "quantity", "created_by", "created_at"})
		n := 2
		for i := 0; i < 2; i++ {
			row.AddRow(i+1, mockData[i].Title, pq.Array(mockData[i].Authors), mockData[i].Publisher, mockData[i].Isbn, mockData[0].Price, mockData[0].Quantity, mockData[i].Created_by, mockCreated_at)
		}

		get := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, title, authors, publisher, isbn, price, quantity, created_by, created_at FROM books ORDER BY id LIMIT $1 OFFSET $2;`))
		get.ExpectQuery().
			WithArgs().
			WillReturnRows(row)

		query := NewDB(db)

		// Act
		results, err := query.SelectAllBooks(params)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, n)
	})

	t.Run("TestSelectAllBooksShouldReturnError", func(t *testing.T) {
		// Arrange
		params := GetAllParams{
			Limit:  1,
			Offset: 2,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		get := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, title, authors, publisher, isbn, price, quantity, created_by, created_at FROM books ORDER BY id LIMIT $1 OFFSET $2;`))
		get.ExpectQuery().
			WithArgs().
			WillReturnError(&pq.Error{Message: "db connection error"})

		query := NewDB(db)

		// Act
		results, err := query.SelectAllBooks(params)

		// Assert
		assert.NotNil(t, err)
		assert.Nil(t, results)
		assert.Len(t, results, 0)
	})
}

func TestSelectBookByID(t *testing.T) {
	t.Run("TestSelectByIDShouldReturnNoError", func(t *testing.T) {
		// Arrange
		mockData := &RequestBook{
			Title:      "mockData",
			Authors:    []string{"Author A", "Author B"},
			Publisher:  "mockPublisher",
			Isbn:       "1234567890",
			Price:      1000,
			Quantity:   10,
			Created_by: "mockAdmin",
		}
		id := uint64(1)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		row := sqlmock.NewRows([]string{"id", "title", "authors", "publisher", "isbn", "price", "quantity", "created_by", "created_at"})
		mockCreated_at := time.Now()
		row.AddRow(1, mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM books WHERE id = $1;`))
		get.ExpectQuery().
			WithArgs(1).
			WillReturnRows(row)

		query := NewDB(db)

		// Act
		result, err := query.SelectBookByID(id)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, mockData.Title, result.Title)
		assert.Equal(t, mockData.Authors, result.Authors)
		assert.Equal(t, mockData.Publisher, result.Publisher)
		assert.Equal(t, mockData.Price, result.Price)
		assert.Equal(t, mockData.Quantity, result.Quantity)
		assert.Equal(t, mockData.Isbn, result.Isbn)
		assert.Equal(t, mockData.Created_by, result.Created_by)
		assert.Equal(t, id, result.Id)
		assert.Equal(t, mockCreated_at, result.Created_at)
	})

	t.Run("TestSelectBookByIDShouldReturnError", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mockData := RequestBook{}
		row := sqlmock.NewRows([]string{"id", "title", "authors", "publisher", "isbn", "price", "quantity", "created_by", "created_at"})

		id := uint64(1)
		mockCreated_at := time.Now()
		row.AddRow(1, mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM books WHERE id = $1;`))
		get.ExpectQuery().
			WithArgs(id).
			WillReturnError(&pq.Error{Message: "db connection error"})

		query := NewDB(db)

		// Act
		result, err := query.SelectBookByID(id)

		// Assert
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdateBook(t *testing.T) {
	t.Run("TestUpdateBookShouldReturnNoError", func(t *testing.T) {
		// Arrange
		id := uint64(1)
		mockData := RequestBook{
			Title:      "mockData",
			Authors:    []string{"Author A", "Author B"},
			Publisher:  "mockPublisher",
			Isbn:       "1234567890",
			Price:      1000,
			Quantity:   10,
			Created_by: "mockAdmin",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mockCreated_at := time.Now()
		row := sqlmock.NewRows([]string{"id", "title", "authors", "publisher", "isbn", "price", "quantity", "created_by", "created_at"})
		row.AddRow(id, mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE books SET title = $1, authors = $2, publisher = $3, isbn = $4, price = $5, quantity = $6, created_by = $7 WHERE id = $8 RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`))
		get.ExpectQuery().
			WithArgs(mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by, id).
			WillReturnRows(row)

		query := NewDB(db)

		// Act
		result, err := query.UpdateBook(id, mockData)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, mockData.Title, result.Title)
		assert.Equal(t, mockData.Authors, result.Authors)
		assert.Equal(t, mockData.Publisher, result.Publisher)
		assert.Equal(t, mockData.Price, result.Price)
		assert.Equal(t, mockData.Quantity, result.Quantity)
		assert.Equal(t, mockData.Isbn, result.Isbn)
		assert.Equal(t, mockData.Created_by, result.Created_by)
		assert.Equal(t, id, result.Id)
		assert.Equal(t, mockCreated_at, result.Created_at)
	})

	t.Run("TestUpdateBookShouldReturnError", func(t *testing.T) {
		// Arrange
		id := uint64(1)
		mockData := RequestBook{
			Title:      "mockData",
			Authors:    []string{"Author A", "Author B"},
			Publisher:  "mockPublisher",
			Isbn:       "1234567890",
			Price:      1000,
			Quantity:   10,
			Created_by: "mockAdmin",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		get := mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE books SET title = $1, authors = $2, publisher = $3, isbn = $4, price = $5, quantity = $6, created_by = $7 WHERE id = $8 RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`))
		get.ExpectQuery().
			WithArgs(id, mockData).
			WillReturnError(&pq.Error{Message: "db connection error"})

		query := NewDB(db)

		// Act
		result, err := query.UpdateBook(id, mockData)

		// Assert j
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestDeleteBook(t *testing.T) {
	t.Run("TestDeleteBookShouldReturnNoError", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		id := uint64(1)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1;`))
		get.ExpectExec().
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		query := NewDB(db)

		// Act
		err = query.DeleteBook(id)

		// Assert
		assert.Nil(t, err)
	})

	t.Run("TestDeleteBookShouldReturnError", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		id := uint64(1)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1;`))
		get.ExpectExec().
			WithArgs(id).
			WillReturnError(&pq.Error{Message: "invalid id"})

		query := NewDB(db)

		// Act
		err = query.DeleteBook(id)

		// Assert
		assert.NotNil(t, err)
	})
}
