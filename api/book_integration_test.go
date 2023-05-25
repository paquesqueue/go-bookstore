//go:build integration

package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const serverPortBook = 80

func setupServerBook(t *testing.T) (*sql.DB, func()) {
	// Setup server
	ec := echo.New()
	db, err := sql.Open("postgres", "postgres://etyaksig:lxqUGmKanXKjhCUYJ-UWqUSzgAZl9g4L@tiny.db.elephantsql.com/etyaksig?sslmode=disable")
	if err != nil {
		assert.Error(t, err)
	}
	go func(e *echo.Echo, db *sql.DB) {
		log := logrus.New()
		storage := NewDB(db)
		service := NewBookService(storage, log)
		handler := NewBookHandlr(service, log)

		e.POST("/books", handler.AddBook)
		e.GET("/books", handler.ListAllBooks)
		e.GET("/books/:id", handler.GetBookByID)
		e.PUT("/books/:id", handler.PutBook)
		e.DELETE("/books/:id", handler.DelBook)
		e.Start(fmt.Sprintf(":%d", serverPortBook))

	}(ec, db)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPortBook), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	return db, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = ec.Shutdown(ctx)
		assert.NoError(t, err)
		err = db.Close()
		assert.NoError(t, err)
	}
}

func TestAddBookIntegration(t *testing.T) {
	_, teardown := setupServerBook(t)
	defer teardown()
	// Arrange
	mockData := RequestBook{
		Title:      "mockDataAd",
		Authors:    []string{"Author A", "Author B"},
		Publisher:  "mockPublisher",
		Isbn:       "1234567890",
		Price:      1000,
		Quantity:   10,
		Created_by: "mockAdmin",
	}
	body, err := json.Marshal(mockData)
	if err != nil {
		assert.Fail(t, "error json marshal")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/books", serverPortBook), strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	result := &ResponseBook{}
	err = json.Unmarshal(byteBody, &result)
	assert.NoError(t, err)

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, int64(0), result.Id)
		assert.Equal(t, mockData.Title, result.Title)
		assert.Equal(t, mockData.Authors, result.Authors)
		assert.Equal(t, mockData.Publisher, result.Publisher)
		assert.Equal(t, mockData.Price, result.Price)
		assert.Equal(t, mockData.Quantity, result.Quantity)
		assert.Equal(t, mockData.Isbn, result.Isbn)
		assert.Equal(t, mockData.Created_by, result.Created_by)
	}
}

func TestSelectBookByIDIntegration(t *testing.T) {
	db, teardown := setupServerBook(t)
	defer teardown()
	// Arrange
	stmt, err := db.Prepare(`INSERT INTO books (title, authors, publisher, isbn, price, quantity, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`)
	assert.NoError(t, err)
	defer stmt.Close()

	mockData := RequestBook{
		Title:      "mockDataGetByID",
		Authors:    []string{"Author A", "Author B"},
		Publisher:  "mockPublisher",
		Isbn:       "1234567890",
		Price:      1000,
		Quantity:   10,
		Created_by: "mockAdmin",
	}
	row := stmt.QueryRow(mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by)

	var response ResponseBook
	err = row.Scan(&response.Id, &response.Title, pq.Array(&response.Authors), &response.Publisher, &response.Isbn, &response.Price, &response.Quantity, &response.Created_by, &response.Created_at)
	assert.NoError(t, err)

	targetUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d/books", serverPortBook))
	assert.NoError(t, err)
	targetUrl = targetUrl.JoinPath(strconv.FormatInt(int64(response.Id), 10))

	req, err := http.NewRequest(http.MethodGet, targetUrl.String(), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	result := &ResponseBook{}
	err = json.Unmarshal(byteBody, &result)
	assert.NoError(t, err)

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, response.Id, result.Id)
		assert.Equal(t, mockData.Title, result.Title)
		assert.Equal(t, mockData.Authors, result.Authors)
		assert.Equal(t, mockData.Publisher, result.Publisher)
		assert.Equal(t, mockData.Price, result.Price)
		assert.Equal(t, mockData.Quantity, result.Quantity)
		assert.Equal(t, mockData.Isbn, result.Isbn)
		assert.Equal(t, mockData.Created_by, result.Created_by)
	}
}

func TestUpdateBookIntegration(t *testing.T) {
	db, teardown := setupServerBook(t)
	defer teardown()
	// Arrange
	stmt, err := db.Prepare(`INSERT INTO books (title, authors, publisher, isbn, price, quantity, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`)
	assert.NoError(t, err)
	defer stmt.Close()

	mockData := RequestBook{
		Title:      "mockData",
		Authors:    []string{"Author A", "Author B"},
		Publisher:  "mockPublisher",
		Isbn:       "1234567890",
		Price:      1000,
		Quantity:   10,
		Created_by: "mockAdmin",
	}
	row := stmt.QueryRow(mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by)

	var response ResponseBook
	err = row.Scan(&response.Id, &response.Title, pq.Array(&response.Authors), &response.Publisher, &response.Isbn, &response.Price, &response.Quantity, &response.Created_by, &response.Created_at)
	assert.NoError(t, err)

	targetUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d/books", serverPortBook))
	assert.NoError(t, err)
	targetUrl = targetUrl.JoinPath(strconv.FormatInt(int64(response.Id), 10))

	reqBody := RequestBook{
		Title:      "newMockData",
		Authors:    []string{"New Author A", "New Author B"},
		Publisher:  "newMockPublisher",
		Isbn:       "01234567890",
		Price:      200,
		Quantity:   1,
		Created_by: "mockAdmin",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		assert.Fail(t, "error json marshal")
	}

	req, err := http.NewRequest(http.MethodPut, targetUrl.String(), strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	result := &ResponseBook{}
	err = json.Unmarshal(byteBody, &result)
	assert.NoError(t, err)

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, response.Id, result.Id)
		assert.Equal(t, reqBody.Title, result.Title)
		assert.Equal(t, reqBody.Authors, result.Authors)
		assert.Equal(t, reqBody.Publisher, result.Publisher)
		assert.Equal(t, reqBody.Price, result.Price)
		assert.Equal(t, reqBody.Quantity, result.Quantity)
		assert.Equal(t, reqBody.Isbn, result.Isbn)
		assert.Equal(t, reqBody.Created_by, result.Created_by)
	}
}

func TestSelectAllBooksIntegration(t *testing.T) {
	db, teardown := setupServerBook(t)
	defer teardown()
	// Arrange
	stmt, err := db.Prepare(`DELETE FROM books;`)
	assert.NoError(t, err)
	defer stmt.Close()

	_, err = stmt.Exec()
	assert.NoError(t, err)

	stmt, err = db.Prepare(`INSERT INTO books (title, authors, publisher, isbn, price, quantity, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7);`)
	assert.NoError(t, err)

	mockData := RequestBook{
		Title:      "mockDataGetAll",
		Authors:    []string{"Author A", "Author B"},
		Publisher:  "mockPublisher",
		Isbn:       "1234567890",
		Price:      1000,
		Quantity:   10,
		Created_by: "mockAdmin",
	}

	_ = stmt.QueryRow(mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by)

	params := RequestGetAll{
		PageId:   1,
		PageSize: 1,
	}
	body, err := json.Marshal(params)
	if err != nil {
		assert.Fail(t, "json marshal error")
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/books", serverPortBook), strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	resp.Body.Close()

	respBody := []ResponseBook{}
	err = json.Unmarshal(byteBody, &respBody)

	assert.NoError(t, err)

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Len(t, respBody, 1)
	}
}

func TestDeleteBooksIntegration(t *testing.T) {
	db, teardown := setupServerBook(t)
	defer teardown()
	// Arrange
	stmt, err := db.Prepare(`INSERT INTO books (title, authors, publisher, isbn, price, quantity, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, title, authors, publisher, isbn, price, quantity, created_by, created_at;`)
	assert.NoError(t, err)
	defer stmt.Close()

	mockData := RequestBook{
		Title:      "mockDataDelete",
		Authors:    []string{"Author A", "Author B"},
		Publisher:  "mockPublisher",
		Isbn:       "1234567890",
		Price:      1000,
		Quantity:   10,
		Created_by: "mockAdmin",
	}
	row := stmt.QueryRow(mockData.Title, pq.Array(mockData.Authors), mockData.Publisher, mockData.Isbn, mockData.Price, mockData.Quantity, mockData.Created_by)
	response := ResponseBook{}
	row.Scan(&response.Id)
	assert.NoError(t, err)

	body, err := json.Marshal(response.Id)
	if err != nil {
		assert.Fail(t, "json marshal error")
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:%d/books", serverPortBook), strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	result := &ResponseBook{}
	err = json.Unmarshal(byteBody, &result)
	assert.NoError(t, err)
	assert.Empty(t, result)
}
