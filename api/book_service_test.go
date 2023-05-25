//go:build unit
package api

import (
	"testing"
	"time"

	c "github.com/paquesqueue/bookstore/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type BookQueriesSuccess struct {
	insertBookCalled  bool
	selectAllBooksCalled bool
	selectBookByIDCalled  bool
	updateBookCalled  bool
	deleteBookCallled bool
}

func (s *BookQueriesSuccess) InsertBook(req RequestBook) (*ResponseBook, error) {
	s.insertBookCalled = true
	resp := &ResponseBook{
		Id:         1,
		Title:      req.Title,
		Authors:    req.Authors,
		Publisher:  req.Publisher,
		Isbn:       req.Isbn,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Created_by: "Admin",
		Created_at: time.Now(),
	}
	return resp, nil
}

func (s *BookQueriesSuccess) SelectAllBooks(params GetAllParams) ([]ResponseBook, error) {
	s.selectAllBooksCalled = true
	resp := []ResponseBook{
		{
			Id:         1,
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
			Created_at: time.Now(),
		},
		{
			Id:         2,
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
			Created_at: time.Now(),
		},
	}
	return resp, nil
}

func (s *BookQueriesSuccess) SelectBookByID(id uint64) (*ResponseBook, error) {
	s.selectBookByIDCalled = true
	resp := &ResponseBook{
		Id:         id,
		Title:      "mockTitle",
		Authors:    []string{"mockAuthors"},
		Publisher:  "mockPublisher",
		Isbn:       "mockIsbn",
		Price:      1000,
		Quantity:   100,
		Created_by: "Admin",
		Created_at: time.Now(),
	}
	return resp, nil
}

func (s *BookQueriesSuccess) UpdateBook(id uint64, req RequestBook) (*ResponseBook, error) {
	s.updateBookCalled = true
	resp := &ResponseBook{
		Id:         id,
		Title:      "mockTitle",
		Authors:    []string{"mockAuthors"},
		Publisher:  "mockPublisher",
		Isbn:       "mockIsbn",
		Price:      1000,
		Quantity:   100,
		Created_by: "Admin",
		Created_at: time.Now(),
	}
	return resp, nil
}

func (s *BookQueriesSuccess) DeleteBook(id uint64) error {
	s.deleteBookCallled = true
	return nil
}

type BookQueriesError struct {
	insertBookCalled  bool
	selectAllBooksCalled bool
	selectBookByIDCalled  bool
	updateBookCalled  bool
	deleteBookCallled bool
}

func (s *BookQueriesError) InsertBook(req RequestBook) (*ResponseBook, error) {
	s.insertBookCalled = true
	return nil, &c.Err{}
}

func (s *BookQueriesError) SelectAllBooks(params GetAllParams) ([]ResponseBook, error) {
	s.selectAllBooksCalled = true
	return nil, &c.Err{}
}

func (s *BookQueriesError) SelectBookByID(id uint64) (*ResponseBook, error) {
	s.selectBookByIDCalled = true
	return nil, &c.Err{}
}

func (s *BookQueriesError) UpdateBook(id uint64, req RequestBook) (*ResponseBook, error) {
	s.updateBookCalled = true
	return nil, &c.Err{}
}

func (s *BookQueriesError) DeleteBook(id uint64) error {
	s.deleteBookCallled = true
	return &c.Err{}
}

func TestAddBook(t *testing.T) {
	t.Run("TestAddBookServiceShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &BookQueriesSuccess{}
		log := logrus.New()
		services := NewBookService(query, log)

		mockData := RequestBook{
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
		}

		// Act
		res, err := services.AddBook(mockData)

		// Assert
		assert.Equal(t, true, query.insertBookCalled)
		assert.NotNil(t, res)
		assert.Nil(t, err)

		assert.NotEqual(t, uint64(0), res.Id)
		assert.Equal(t, mockData.Title, res.Title)
		assert.Equal(t, mockData.Authors, res.Authors)
		assert.Equal(t, mockData.Publisher, res.Publisher)
		assert.Equal(t, mockData.Isbn, res.Isbn)
		assert.Equal(t, mockData.Price, res.Price)
		assert.Equal(t, mockData.Quantity, res.Quantity)
		assert.Equal(t, mockData.Created_by, res.Created_by)
		assert.NotEmpty(t, res.Created_at)
	})

	t.Run("TestAddBookServiceShouldReturnError", func(t *testing.T) {
		// Arrange
		query := &BookQueriesError{}
		log := logrus.New()
		services := NewBookService(query, log)

		mockData := RequestBook{}

		// Act
		res, err := services.AddBook(mockData)

		// Assert
		assert.Equal(t, true, query.insertBookCalled)
		assert.Nil(t, res)
		assert.NotNil(t, err)

	})
}

func TestListAllBooks(t *testing.T) {
	t.Run("TestListAllBooksServiceShouldReturnNoError", func(t *testing.T) {
		// Arrange
		params := GetAllParams{
			Limit:  1,
			Offset: 2,
		}
		query := &BookQueriesSuccess{}
		log := logrus.New()
		services := NewBookService(query, log)

		// Act
		res, err := services.ListAllBooks(params)

		// Assert
		assert.Equal(t, true, query.selectAllBooksCalled)
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("TestListAllBooksServiceShouldReturnError", func(t *testing.T) {
		// Arrange
		params := GetAllParams{
			Limit:  1,
			Offset: 2,
		}
		query := &BookQueriesError{}
		log := logrus.New()
		services := NewBookService(query, log)

		// Act
		res, err := services.ListAllBooks(params)

		// Assert
		assert.Equal(t, true, query.selectAllBooksCalled)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestGetBookByID(t *testing.T) {
	t.Run("TestGetBookByIDShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &BookQueriesSuccess{}
		log := logrus.New()
		services := NewBookService(query, log)

		id := uint64(1)

		// Act
		res, err := services.GetBookByID(id)

		// Assert
		assert.Equal(t, true, query.selectBookByIDCalled)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id, res.Id)
	})

	t.Run("TestGetBookByIDShouldReturnError", func(t *testing.T) {
		// Arrange
		query := &BookQueriesError{}
		log := logrus.New()
		services := NewBookService(query, log)

		id := uint64(0)

		// Act
		res, err := services.GetBookByID(id)

		// Assert
		assert.Equal(t, true, query.selectBookByIDCalled)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestPutBook(t *testing.T) {
	t.Run("TestUpdateBookShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &BookQueriesSuccess{}
		log := logrus.New()
		services := NewBookService(query, log)

		id := uint64(1)
		mockData := RequestBook{
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
		}

		// Act
		res, err := services.PutBook(id, mockData)

		// Assert
		assert.Equal(t, true, query.updateBookCalled)
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Equal(t, id, res.Id)
		assert.Equal(t, mockData.Title, res.Title)
		assert.Equal(t, mockData.Authors, res.Authors)
		assert.Equal(t, mockData.Publisher, res.Publisher)
		assert.Equal(t, mockData.Isbn, res.Isbn)
		assert.Equal(t, mockData.Price, res.Price)
		assert.Equal(t, mockData.Quantity, res.Quantity)
		assert.Equal(t, mockData.Created_by, res.Created_by)
		assert.NotEmpty(t, res.Created_at)
	})

	t.Run("TestUpdateBookShouldReturnError", func(T *testing.T) {
		// Arrange
		query := &BookQueriesError{}
		log := logrus.New()
		services := NewBookService(query, log)

		id := uint64(1)
		mockData := RequestBook{
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
		}

		// Act
		res, err := services.PutBook(id, mockData)

		// Assert
		assert.Equal(t, true, query.updateBookCalled)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestDelBook(t *testing.T) {
	t.Run("TestDelBookShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &BookQueriesSuccess{}
		log := logrus.New()
		services := NewBookService(query, log)

		id := uint64(1)

		// Act
		err := services.DelBook(id)

		// Assert
		assert.Equal(t, true, query.deleteBookCallled)
		assert.NoError(t, err)
	})

	t.Run("TestDelBookShouldReturnError", func(t *testing.T) {
		// Arrange
		query := &BookQueriesError{}
		log := logrus.New()

		services := NewBookService(query, log)

		id := uint64(1)

		// Act
		err := services.DelBook(id)

		// Assert
		assert.Equal(t, true, query.deleteBookCallled)
		assert.Error(t, err)
	})
}
