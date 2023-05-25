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
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/paquesqueue/bookstore/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const serverPortUser = 80

func setupServerUser(t *testing.T) (*sql.DB, func()) {
	// Setup server
	ec := echo.New()
	db, err := sql.Open("postgres", "postgres://etyaksig:lxqUGmKanXKjhCUYJ-UWqUSzgAZl9g4L@tiny.db.elephantsql.com/etyaksig?sslmode=disable")
	if err != nil {
		assert.Error(t, err)
	}
	go func(e *echo.Echo, db *sql.DB) {
		log := logrus.New()
		storage := NewDB(db)
		service := NewUserService(storage, log)
		handler := NewUserHandler(service, log)

		e.POST("/users", handler.AddUser)
		e.GET("/users/:username", handler.GetUser)
		e.PUT("/users/:username", handler.PutUser)
		e.DELETE("/users/:username", handler.DeleteUser)
		e.Start(fmt.Sprintf(":%d", serverPortUser))

	}(ec, db)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPortUser), 30*time.Second)
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

func TestAddUserIntegration(t *testing.T) {
	_, teardown := setupServerUser(t)
	defer teardown()
	// Arrange

	mockData := RequestUser{
		Username: utils.RandomUsername(),
		Password: utils.RandomPassword(),
		Email:    utils.RandomEmail(),
		Fullname: utils.RandomFullname(),
	}
	body, err := json.Marshal(mockData)
	if err != nil {
		assert.Fail(t, "error json marshal")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/users", serverPortUser), strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	result := &ResponseUser{}
	err = json.Unmarshal(byteBody, &result)
	assert.NoError(t, err)

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, mockData.Username, result.Username)
		assert.Equal(t, mockData.Email, result.Email)
		assert.Equal(t, mockData.Fullname, result.Fullname)
		assert.NotEqual(t, mockData.Password, result.HashedPassword)
		assert.NotEmpty(t, result.CreatedAt)
	}
}

func TestSelectUserIntegration(t *testing.T) {
	db, teardown := setupServerUser(t)
	defer teardown()
	// Arrange
	stmt, err := db.Prepare(`DELETE FROM users;`)
	assert.NoError(t, err)
	defer stmt.Close()

	_, err = stmt.Exec()
	assert.NoError(t, err)

	stmt, err = db.Prepare(`INSERT INTO users (username, email, fullname, hashed_password) VALUES ($1, $2, $3, $4) RETURNING username, email, fullname, hashed_password, created_at;`)
	assert.NoError(t, err)

	mockData := RequestUser{
		Username: utils.RandomUsername(),
		Password: utils.RandomPassword(),
		Email:    utils.RandomEmail(),
		Fullname: utils.RandomFullname(),
	}

	row := stmt.QueryRow(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password)

	var response ResponseUser
	err = row.Scan(&response.Username, &response.Email, &response.Fullname, &response.HashedPassword, &response.CreatedAt)
	assert.NoError(t, err)

	targetUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d/users", serverPortUser))
	assert.NoError(t, err)
	targetUrl = targetUrl.JoinPath(response.Username)

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

	result := &ResponseUser{}
	err = json.Unmarshal(byteBody, &result)
	assert.NoError(t, err)

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, mockData.Username, result.Username)
		assert.Equal(t, mockData.Email, result.Email)
		assert.Equal(t, mockData.Fullname, result.Fullname)
		assert.Equal(t, mockData.Password, result.HashedPassword)
		assert.NotEmpty(t, result.CreatedAt)
	}
}

func TestUpdateUserIntegration(t *testing.T) {
	db, teardown := setupServerUser(t)
	defer teardown()
	stmt, err := db.Prepare(`DELETE FROM users;`)
	assert.NoError(t, err)
	defer stmt.Close()

	_, err = stmt.Exec()
	assert.NoError(t, err)

	// Arrange
	stmt, err = db.Prepare(`INSERT INTO users (username, email, fullname, hashed_password) VALUES ($1, $2, $3, $4) RETURNING username, email, fullname, hashed_password, created_at;`)
	assert.NoError(t, err)

	mockData := RequestUser{
		Username: utils.RandomUsername(),
		Password: utils.RandomPassword(),
		Email:    utils.RandomEmail(),
		Fullname: utils.RandomFullname(),
	}

	row := stmt.QueryRow(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password)

	var response ResponseUser
	err = row.Scan(&response.Username, &response.Email, &response.Fullname, &response.HashedPassword, &response.CreatedAt)
	assert.NoError(t, err)

	targetUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d/users", serverPortUser))
	assert.NoError(t, err)
	targetUrl = targetUrl.JoinPath(response.Username)

	reqBody := RequestUser{
		Username: mockData.Username,
		Password: utils.RandomPassword(),
		Email:    utils.RandomEmail(),
		Fullname: utils.RandomFullname(),
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

	result := &ResponseUser{}
	err = json.Unmarshal(byteBody, &result)
	assert.NoError(t, err)

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, reqBody.Username, result.Username)
		assert.Equal(t, reqBody.Email, result.Email)
		assert.Equal(t, reqBody.Fullname, result.Fullname)
		assert.NotEqual(t, reqBody.Password, result.HashedPassword)
		assert.NotEmpty(t, result.CreatedAt)
	}
}
