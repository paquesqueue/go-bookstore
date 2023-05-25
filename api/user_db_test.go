//go:build unit

package api

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/paquesqueue/bookstore/utils"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	t.Run("TestInsertUserShouldReturnNoError", func(t *testing.T) {
		// Arrange
		hashedPassword, err := utils.HashPassword("123456")
		assert.NoError(t, err)

		mockData := RequestUser{
			Username: "tester",
			Password: hashedPassword,
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		row := sqlmock.NewRows([]string{"username", "email", "fullname", "hashed_password", "created_at"})
		mockCreated_at := time.Now()
		row.AddRow(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO users (username, email, fullname, hashed_password) VALUES ($1, $2, $3, $4) RETURNING username, email, fullname, hashed_password, created_at;`))
		get.ExpectQuery().
			WithArgs(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password).
			WillReturnRows(row)

		query := NewDB(db)

		// Act
		resp, err := query.InsertUser(mockData)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, mockData.Username, resp.Username)
			assert.Equal(t, mockData.Email, resp.Email)
			assert.Equal(t, mockData.Fullname, resp.Fullname)
			assert.Equal(t, mockData.Password, resp.HashedPassword)
			assert.NotEmpty(t, resp.CreatedAt)
		}
	})

	t.Run("TestInsertUserShouldReturnError", func(t *testing.T) {
		// Arrange
		mockData := RequestUser{}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		get := mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO users (username, email, fullname, hashed_password) VALUES ($1, $2, $3, $4) RETURNING username, email, fullname, hashed_password, created_at;`))
		get.ExpectQuery().
			WithArgs(mockData).
			WillReturnError(&pq.Error{Message: "db connection error"})

		query := NewDB(db)

		// Act
		resp, err := query.InsertUser(mockData)

		// Assert
		if assert.NotNil(t, err) {
			assert.Empty(t, resp)
		}
	})
}

func TestSelectUser(t *testing.T) {
	t.Run("TestSelectUserShouldReturnNoError", func(t *testing.T) {
		// Arrange
		hashedPassword, err := utils.HashPassword("123456")
		assert.NoError(t, err)

		mockData := RequestUser{
			Username: "tester",
			Password: hashedPassword,
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		row := sqlmock.NewRows([]string{"username", "email", "fullname", "hashed_password", "created_at"})
		mockCreated_at := time.Now()
		row.AddRow(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM users WHERE username = $1;`))
		get.ExpectQuery().
			WithArgs(mockData.Username).
			WillReturnRows(row)

		query := NewDB(db)

		// Act
		resp, err := query.SelectUser(mockData.Username)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, mockData.Username, resp.Username)
			assert.Equal(t, mockData.Email, resp.Email)
			assert.Equal(t, mockData.Fullname, resp.Fullname)
			assert.Equal(t, mockData.Password, resp.HashedPassword)
			assert.NotEmpty(t, resp.CreatedAt)
		}
	})

	t.Run("TestSelectUserShouldReturnError", func(t *testing.T) {
		// Arrange
		hashedPassword, err := utils.HashPassword("123456")
		assert.NoError(t, err)

		mockData := RequestUser{
			Username: "tester",
			Password: hashedPassword,
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		get := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM users WHERE username = $1;`))
		get.ExpectQuery().
			WithArgs(mockData.Username).
			WillReturnError(&pq.Error{Message: "not found error"})

		query := NewDB(db)

		// Act
		resp, err := query.SelectUser(mockData.Username)

		// Assert
		if assert.NotNil(t, err) {
			assert.Empty(t, resp)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("TestUpdatetUserShouldReturnNoError", func(t *testing.T) {
		// Arrange
		hashedPassword, err := utils.HashPassword("123456")
		assert.NoError(t, err)

		username := "tester"

		mockData := RequestUser{
			Username: "tester",
			Password: hashedPassword,
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		row := sqlmock.NewRows([]string{"username", "email", "fullname", "hashed_password", "created_at"})
		mockCreated_at := time.Now()
		row.AddRow(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET username = $1, email = $2, fullname = $3, hashed_password = $4 WHERE username = $5 RETURNING username, email, fullname, hashed_password, created_at;`))
		get.ExpectQuery().
			WithArgs(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password, mockData.Username).
			WillReturnRows(row)

		query := NewDB(db)

		// Act
		resp, err := query.UpdateUser(username, mockData)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, mockData.Username, resp.Username)
			assert.Equal(t, mockData.Email, resp.Email)
			assert.Equal(t, mockData.Fullname, resp.Fullname)
			assert.Equal(t, mockData.Password, resp.HashedPassword)
			assert.NotEmpty(t, resp.CreatedAt)
		}
	})

	t.Run("TestUpdatetUserShouldReturnError", func(t *testing.T) {
		// Arrange
		mockData := RequestUser{}
		username := "tester"
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		row := sqlmock.NewRows([]string{"username", "email", "fullname", "hashed_password", "created_at"})
		mockCreated_at := time.Now()
		row.AddRow(mockData.Username, mockData.Email, mockData.Fullname, mockData.Password, mockCreated_at)

		get := mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET username = $1, email = $2, fullname = $3, hashed_password = $4 WHERE username = $5 RETURNING username, email, fullname, hashed_password, created_at;`))
		get.ExpectQuery().
			WithArgs(mockData.Email, mockData.Fullname, mockData.Password, mockData.Username).
			WillReturnError(&pq.Error{Message: "db connection error"})

		query := NewDB(db)

		// Act
		resp, err := query.UpdateUser(username, mockData)

		// Assert
		if assert.NotNil(t, err) {
			assert.Empty(t, resp)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("TestDeleteUserShouldReturnNoError", func(t *testing.T) {
		// Arrange
		mockUsername := "tester"

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		get := mock.ExpectPrepare(regexp.QuoteMeta(`DELETE FROM users WHERE username = $1`))
		get.ExpectExec().
			WithArgs(mockUsername).
			WillReturnResult(sqlmock.NewResult(1, 1))

		query := NewDB(db)

		// Act
		err = query.DeleteUser(mockUsername)

		// Assert
		assert.Nil(t, err)
	})

	t.Run("TestDeleteUserShouldReturnError", func(t *testing.T) {
		// Arrange
		mockUsername := "tester"

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		get := mock.ExpectPrepare(regexp.QuoteMeta(`DELETE FROM users WHERE username = $1`))
		get.ExpectExec().
			WithArgs(mockUsername).
			WillReturnError(&pq.Error{Message: "db connection error"})

		query := NewDB(db)

		// Act
		err = query.DeleteUser(mockUsername)

		// Assert
		assert.NotNil(t, err)
	})
}
