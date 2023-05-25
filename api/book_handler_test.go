//go:build unit
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	c "github.com/paquesqueue/bookstore/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type BookHandlrSuccess struct {
	addBookCalled   bool
	getBookByIDCalled   bool
	listAllBooksCalled bool
	putBookCalled   bool
	delBookCalled   bool
}

func (h *BookHandlrSuccess) AddBook(req RequestBook) (*ResponseBook, error) {
	h.addBookCalled = true
	res := &ResponseBook{
		Id:         uint64(1),
		Title:      req.Title,
		Publisher:  req.Publisher,
		Isbn:       req.Isbn,
		Authors:    req.Authors,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Created_by: "Admin",
		Created_at: time.Now(),
	}
	return res, nil
}

func (h *BookHandlrSuccess) GetBookByID(id uint64) (*ResponseBook, error) {
	h.getBookByIDCalled = true
	res := &ResponseBook{
		Id:         1,
		Title:      "mockTitle",
		Authors:    []string{"mockAuthors"},
		Publisher:  "mockPublisher",
		Isbn:       "mockIsbn",
		Price:      1000,
		Quantity:   100,
		Created_by: "Admin",
		Created_at: time.Now(),
	}
	return res, nil
}

func (h *BookHandlrSuccess) ListAllBooks(params GetAllParams) ([]ResponseBook, error) {
	h.listAllBooksCalled = true
	res := []ResponseBook{
		{Id: 1,
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
	return res, nil
}

func (h *BookHandlrSuccess) PutBook(id uint64, req RequestBook) (*ResponseBook, error) {
	h.putBookCalled = true
	res := &ResponseBook{
		Id:         id,
		Title:      req.Title,
		Publisher:  req.Publisher,
		Isbn:       req.Isbn,
		Authors:    req.Authors,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Created_by: "Admin",
		Created_at: time.Now(),
	}
	return res, nil
}

func (h *BookHandlrSuccess) DelBook(id uint64) error {
	h.delBookCalled = true
	return nil
}

type BookHandlrError struct {
	addBookCalled   bool
	getBookByIDCalled   bool
	listAllBooksCalled bool
	putBookCalled   bool
	delBookCalled   bool
	statusCodeError int
}

func (h *BookHandlrError) AddBook(req RequestBook) (*ResponseBook, error) {
	h.addBookCalled = true
	return nil, &c.Err{Code: h.statusCodeError}
}

func (h *BookHandlrError) GetBookByID(id uint64) (*ResponseBook, error) {
	h.getBookByIDCalled = true
	return nil, &c.Err{Code: h.statusCodeError}
}

func (h *BookHandlrError) ListAllBooks(params GetAllParams) ([]ResponseBook, error) {
	h.listAllBooksCalled = true
	return nil, &c.Err{Code: h.statusCodeError}
}

func (h *BookHandlrError) PutBook(id uint64, req RequestBook) (*ResponseBook, error) {
	h.putBookCalled = true
	return nil, &c.Err{Code: h.statusCodeError}
}

func (h *BookHandlrError) DelBook(uid uint64) error {
	h.delBookCalled = true
	return &c.Err{Code: h.statusCodeError}
}

func TestAddBookHandler(t *testing.T) {
	t.Run("TestAddBookHandlerShouldReturnHTTPStatus201", func(t *testing.T) {
		// Arrange
		reqBody := RequestBook{
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal error")
		}

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)

		handlrServ := &BookHandlrSuccess{}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err = handler.AddBook(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			res := &ResponseBook{}
			json.Unmarshal(rec.Body.Bytes(), res)

			assert.NotEqual(t, uint64(0), res.Id)
			assert.Equal(t, reqBody.Title, res.Title)
			assert.Equal(t, reqBody.Authors, res.Authors)
			assert.Equal(t, reqBody.Publisher, res.Publisher)
			assert.Equal(t, reqBody.Isbn, res.Isbn)
			assert.Equal(t, reqBody.Price, res.Price)
			assert.Equal(t, reqBody.Quantity, res.Quantity)
			assert.Equal(t, reqBody.Created_by, res.Created_by)
			assert.NotEmpty(t, res.Created_at)
		}
	})

	t.Run("TestAddBookHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		reqBody := &RequestBook{}
		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal fail")
		}

		req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handlrServ := &BookHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err = handler.AddBook(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, handlrServ.statusCodeError, rec.Code)
		}
	})
}

func TestListAllBooksHandler(t *testing.T) {
	t.Run("TestListAllBooksHandlerShouldReturnHTTPStatus200", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handlrServ := &BookHandlrSuccess{}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err := handler.ListAllBooks(ctx)

		// Assert
		if assert.NoError(t, err) {
			res := []ResponseBook{}
			json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Len(t, res, 2)
		}
	})

	t.Run("TestListAllBooksHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handlrServ := &BookHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err := handler.ListAllBooks(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, handlrServ.statusCodeError, rec.Code)
		}
	})
}

func TestGetBookByIDHandler(t *testing.T) {
	t.Run("TestGetBookByIDHandlerShouldReturnHTTPStatus200", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handlrServ := &BookHandlrSuccess{}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		id := int64(1)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(strconv.FormatInt(id, 10))

		// Act
		err := handler.GetBookByID(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)

			res := &ResponseBook{}
			json.Unmarshal(rec.Body.Bytes(), &res)

			assert.Equal(t, uint64(id), res.Id)
			assert.NotEmpty(t, res.Title)
			assert.NotEmpty(t, res.Authors)
			assert.NotEmpty(t, res.Publisher)
			assert.NotEmpty(t, res.Isbn)
			assert.NotEmpty(t, res.Price)
			assert.NotEmpty(t, res.Quantity)
			assert.NotEmpty(t, res.Created_by)
		}
	})

	t.Run("TestGetBookByIDHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		id := int64(1)
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(strconv.FormatInt(id, 10))

		handlrServ := &BookHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err := handler.GetBookByID(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, handlrServ.statusCodeError, rec.Code)
		}
	})
}

func TestPutBookHandler(t *testing.T) {
	t.Run("TestPutBookHandlerShouldReturnHTTPStatus200", func(t *testing.T) {
		// Act
		reqBody := RequestBook{
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal fail")
		}

		req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		id := int64(1)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(strconv.FormatInt(id, 10))

		handlrServ := &BookHandlrSuccess{}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err = handler.PutBook(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)

			res := &ResponseBook{}
			json.Unmarshal(rec.Body.Bytes(), &res)

			assert.Equal(t, uint64(id), res.Id)
			assert.NotEmpty(t, res.Title)
			assert.NotEmpty(t, res.Authors)
			assert.NotEmpty(t, res.Publisher)
			assert.NotEmpty(t, res.Isbn)
			assert.NotEmpty(t, res.Price)
			assert.NotEmpty(t, res.Quantity)
			assert.NotEmpty(t, res.Created_by)
		}
	})

	t.Run("TestPutBookHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		reqBody := RequestBook{
			Title:      "mockTitle",
			Authors:    []string{"mockAuthors"},
			Publisher:  "mockPublisher",
			Isbn:       "mockIsbn",
			Price:      1000,
			Quantity:   100,
			Created_by: "Admin",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal fail")
		}

		req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		id := int64(1)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(strconv.FormatInt(id, 10))

		handlrServ := &BookHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err = handler.PutBook(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, handlrServ.statusCodeError, rec.Code)
		}
	})

	t.Run("TestPutBookHandlerShouldReturnHTTPStatus400", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/books", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)

		handlrServ := &BookHandlrError{statusCodeError: http.StatusBadRequest}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err := handler.PutBook(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestDelBookHandler(t *testing.T) {
	t.Run("TestDelBookHandlerShouldReturnHTTPStatus200", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodDelete, "/books/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		id := int64(1)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(strconv.FormatInt(id, 10))

		handlrServ := &BookHandlrSuccess{}
		log := logrus.New()

		handler := NewBookHandlr(handlrServ, log)
		err := handler.DelBook(ctx)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("TestDelBookHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodDelete, "/books", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		id := int64(1)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(strconv.FormatInt(id, 10))

		handlrServ := &BookHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err := handler.DelBook(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("TestDelBookHandlerShouldReturnHTTPStatus400", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodDelete, "/books", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handlrServ := &BookHandlrError{statusCodeError: http.StatusBadRequest}
		log := logrus.New()
		handler := NewBookHandlr(handlrServ, log)

		// Act
		err := handler.DelBook(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}
