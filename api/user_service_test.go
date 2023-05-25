//go:build unit

package api

import (
	"testing"
	"time"

	c "github.com/paquesqueue/bookstore/common"
	"github.com/paquesqueue/bookstore/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type UserQueriesSuccess struct {
	insertUserCalled bool
	selectUserCalled bool
	updateUserCalled bool
	deleteUserCalled bool
}

func (s *UserQueriesSuccess) InsertUser(req RequestUser) (ResponseUser, error) {
	s.insertUserCalled = true
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return ResponseUser{}, nil
	}
	return ResponseUser{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
		Fullname:       req.Fullname,

		CreatedAt: time.Now(),
	}, nil
}

func (s *UserQueriesSuccess) SelectUser(username string) (ResponseUser, error) {
	s.selectUserCalled = true
	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		return ResponseUser{}, nil
	}
	return ResponseUser{
		Username:       "tester",
		HashedPassword: hashedPassword,
		Email:          "tester@email.com",
		Fullname:       "tester testing",

		CreatedAt: time.Now(),
	}, nil
}

func (s *UserQueriesSuccess) UpdateUser(username string, req RequestUser) (ResponseUser, error) {
	s.updateUserCalled = true
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return ResponseUser{}, nil
	}
	return ResponseUser{
		Username:       "tester",
		HashedPassword: hashedPassword,
		Email:          "tester@email.com",
		Fullname:       "tester testing",

		CreatedAt: time.Now(),
	}, nil
}

func (s *UserQueriesSuccess) DeleteUser(username string) error {
	s.deleteUserCalled = true
	return nil
}

type UserQueriesError struct {
	insertUserCalled bool
	selectUserCalled bool
	updateUserCalled bool
	deleteUserCalled bool
}

func (s *UserQueriesError) InsertUser(req RequestUser) (ResponseUser, error) {
	s.insertUserCalled = true
	return ResponseUser{}, &c.Err{}
}

func (s *UserQueriesError) SelectUser(username string) (ResponseUser, error) {
	s.selectUserCalled = true
	return ResponseUser{}, &c.Err{}
}

func (s *UserQueriesError) UpdateUser(username string, req RequestUser) (ResponseUser, error) {
	s.updateUserCalled = true
	return ResponseUser{}, &c.Err{}
}

func (s *UserQueriesError) DeleteUser(username string) error {
	s.deleteUserCalled = true
	return &c.Err{}
}

func TestAddUser(t *testing.T) {
	t.Run("TestAddUserServiceShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesSuccess{}
		log := logrus.New()

		services := NewUserService(query, log)

		mockData := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		// Act
		resp, err := services.AddUser(mockData)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, true, query.insertUserCalled)
			assert.Equal(t, mockData.Username, resp.Username)
			assert.Equal(t, mockData.Email, resp.Email)
			assert.Equal(t, mockData.Fullname, resp.Fullname)
			assert.NotEqual(t, mockData.Password, resp.HashedPassword)

			assert.NotEmpty(t, resp.CreatedAt)
		}
	})
	t.Run("TestAddUserServiceShouldReturnError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesError{}
		log := logrus.New()

		services := NewUserService(query, log)

		mockData := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		// Act
		resp, err := services.AddUser(mockData)

		// Assert
		if assert.NotNil(t, err) {
			assert.Equal(t, true, query.insertUserCalled)
			assert.Empty(t, resp)
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("TestGetUserServiceShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesSuccess{}
		log := logrus.New()

		services := NewUserService(query, log)

		mockData := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		// Act
		resp, err := services.GetUser(mockData.Username)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, true, query.selectUserCalled)
			assert.Equal(t, mockData.Username, resp.Username)
			assert.Equal(t, mockData.Email, resp.Email)
			assert.Equal(t, mockData.Fullname, resp.Fullname)
			assert.NotEqual(t, mockData.Password, resp.HashedPassword)

			assert.NotEmpty(t, resp.CreatedAt)
		}
	})
	t.Run("TestGetUserServiceShouldReturnError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesError{}
		log := logrus.New()

		services := NewUserService(query, log)

		mockUsername := "tester"

		// Act
		resp, err := services.GetUser(mockUsername)

		// Assert
		if assert.NotNil(t, err) {
			assert.Equal(t, true, query.selectUserCalled)
			assert.Empty(t, resp)
		}
	})
}

func TestPutUser(t *testing.T) {
	t.Run("TestPutUserServiceShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesSuccess{}
		log := logrus.New()

		services := NewUserService(query, log)
		username := "tester"
		mockData := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		// Act
		resp, err := services.PutUser(username, mockData)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, true, query.updateUserCalled)
			assert.Equal(t, mockData.Username, resp.Username)
			assert.Equal(t, mockData.Email, resp.Email)
			assert.Equal(t, mockData.Fullname, resp.Fullname)
			assert.NotEqual(t, mockData.Password, resp.HashedPassword)

			assert.NotEmpty(t, resp.CreatedAt)
		}
	})
	t.Run("TestPutUserServiceShouldReturnError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesError{}
		log := logrus.New()

		services := NewUserService(query, log)
		username := "tester"
		mockData := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		// Act
		resp, err := services.PutUser(username, mockData)

		// Assert
		if assert.NotNil(t, err) {
			assert.Equal(t, true, query.updateUserCalled)
			assert.Empty(t, resp)
		}
	})
}

func TestDelUser(t *testing.T) {
	t.Run("TestDelUserServiceShouldReturnNoError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesSuccess{}
		log := logrus.New()

		services := NewUserService(query, log)
		mockUsername := "tester"

		// Act
		err := services.DeleteUser(mockUsername)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, true, query.deleteUserCalled)
		}
	})
	t.Run("TestDelUserServiceShouldReturnError", func(t *testing.T) {
		// Arrange
		query := &UserQueriesError{}
		log := logrus.New()

		services := NewUserService(query, log)

		mockUsername := "tester"

		// Act
		err := services.DeleteUser(mockUsername)

		// Assert
		if assert.NotNil(t, err) {
			assert.Equal(t, true, query.deleteUserCalled)
		}
	})
}
