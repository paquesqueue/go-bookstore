//go:build unit

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	c "github.com/paquesqueue/bookstore/common"
	"github.com/paquesqueue/bookstore/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type UserHandlrSuccess struct {
	addUserCalled bool
	getUserCalled bool
	putUserCalled bool
	delUserCalled bool
}

func (s *UserHandlrSuccess) AddUser(req RequestUser) (ResponseUser, error) {
	s.addUserCalled = true

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

func (s *UserHandlrSuccess) GetUser(username string) (ResponseUser, error) {
	s.getUserCalled = true

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

func (s *UserHandlrSuccess) PutUser(username string, req RequestUser) (ResponseUser, error) {
	s.putUserCalled = true
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

func (s *UserHandlrSuccess) DeleteUser(username string) error {
	s.delUserCalled = true
	return nil
}

type UserHandlrError struct {
	addUserCalled   bool
	getUserCalled   bool
	putUserCalled   bool
	delUserCalled   bool
	statusCodeError int
}

func (s *UserHandlrError) AddUser(req RequestUser) (ResponseUser, error) {
	s.addUserCalled = true
	return ResponseUser{}, &c.Err{Code: s.statusCodeError}
}

func (s *UserHandlrError) GetUser(username string) (ResponseUser, error) {
	s.getUserCalled = true
	return ResponseUser{}, &c.Err{Code: s.statusCodeError}
}

func (s *UserHandlrError) PutUser(username string, req RequestUser) (ResponseUser, error) {
	s.putUserCalled = true
	return ResponseUser{}, &c.Err{Code: s.statusCodeError}
}

func (s *UserHandlrError) DeleteUser(username string) error {
	s.delUserCalled = true
	return &c.Err{Code: s.statusCodeError}
}

func TestAddUserHandler(t *testing.T) {
	t.Run("TestAddUserHandlerShouldReturnHTTPStatus201", func(t *testing.T) {
		// Arrange
		reqBody := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal error")
		}
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		handlrServ := &UserHandlrSuccess{}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		e := echo.New()
		ctx := e.NewContext(req, rec)

		// Act
		err = handler.AddUser(ctx)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, true, handlrServ.addUserCalled)

			resp := ResponseUser{}
			json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.Equal(t, reqBody.Username, resp.Username)
			assert.Equal(t, reqBody.Email, resp.Email)
			assert.NotEqual(t, reqBody.Password, resp.HashedPassword)
			assert.Equal(t, reqBody.Fullname, resp.Fullname)

			assert.NotEmpty(t, resp.CreatedAt)
		}
	})

	t.Run("TestAddUserHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		reqBody := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal error")
		}

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		handlrServ := &UserHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()

		handler := NewUserHandler(handlrServ, log)

		e := echo.New()
		ctx := e.NewContext(req, rec)

		// Act
		err = handler.AddUser(ctx)

		// Assert
		if assert.Nil(t, err) {
			assert.Equal(t, true, handlrServ.addUserCalled)
			assert.Equal(t, http.StatusInternalServerError, rec.Code)

			resp := ResponseUser{}
			json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.Empty(t, resp)
		}
	})
}

func TestGetUserHandler(t *testing.T) {
	t.Run("TestGetUserHandlerShouldReturnHTTPStatus200", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handlrServ := &UserHandlrSuccess{}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		username := "tester"
		ctx.SetPath("/:username")
		ctx.SetParamNames("username")
		ctx.SetParamValues(username)

		// Act
		err := handler.GetUser(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, true, handlrServ.getUserCalled)
			assert.Equal(t, http.StatusOK, rec.Code)

			resp := &ResponseUser{}
			json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.NotEmpty(t, resp.Username)
			assert.NotEmpty(t, resp.Email)
			assert.NotEmpty(t, resp.HashedPassword)
			assert.NotEmpty(t, resp.Fullname)

			assert.NotEmpty(t, resp.CreatedAt)
		}
	})

	t.Run("TestGetUserHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)
		username := "tester"
		ctx.SetPath("/:username")
		ctx.SetParamNames("username")
		ctx.SetParamValues(username)

		handlrServ := &UserHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		// Act
		err := handler.GetUser(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, handlrServ.statusCodeError, rec.Code)
		}
	})
}

func TestPutUserHandler(t *testing.T) {
	t.Run("TestPutUserhHandlerShouldReturnHTTPStatus200", func(t *testing.T) {
		// Act
		reqBody := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal fail")
		}

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		username := "tester"
		ctx.SetPath("/:username")
		ctx.SetParamNames("username")
		ctx.SetParamValues(username)

		handlrServ := &UserHandlrSuccess{}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		// Act
		err = handler.PutUser(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)

			resp := &ResponseUser{}
			json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.Equal(t, reqBody.Username, resp.Username)
			assert.Equal(t, reqBody.Email, resp.Email)
			assert.NotEqual(t, reqBody.Password, resp.HashedPassword)
			assert.Equal(t, reqBody.Fullname, resp.Fullname)

			assert.NotEmpty(t, resp.CreatedAt)
		}
	})

	t.Run("TestPutUserHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		reqBody := RequestUser{
			Username: "tester",
			Password: "123456",
			Email:    "tester@email.com",
			Fullname: "tester testing",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal fail")
		}

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		username := "tester"
		ctx.SetPath("/:username")
		ctx.SetParamNames("username")
		ctx.SetParamValues(username)

		handlrServ := &UserHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		// Act
		err = handler.PutUser(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, handlrServ.statusCodeError, rec.Code)
		}
	})

	t.Run("TestPutUserHandlerShouldReturnHTTPStatus400", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		ctx := e.NewContext(req, rec)

		handlrServ := &UserHandlrError{statusCodeError: http.StatusBadRequest}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		// Act
		err := handler.PutUser(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestDeleteUserHandler(t *testing.T) {
	t.Run("TestDeleteUserHandlerShouldReturnHTTPStatus200", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodDelete, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		username := "tester"
		ctx.SetPath("/:username")
		ctx.SetParamNames("username")
		ctx.SetParamValues(username)

		handlrServ := &UserHandlrSuccess{}
		log := logrus.New()

		handler := NewUserHandler(handlrServ, log)
		err := handler.DeleteUser(ctx)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("TestDeleteUserHandlerShouldReturnHTTPStatus500", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodDelete, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		username := "tester"
		ctx.SetPath("/:username")
		ctx.SetParamNames("username")
		ctx.SetParamValues(username)

		handlrServ := &UserHandlrError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		// Act
		err := handler.DeleteUser(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("TestDelBookHandlerShouldReturnHTTPStatus400", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodDelete, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handlrServ := &UserHandlrError{statusCodeError: http.StatusBadRequest}
		log := logrus.New()
		handler := NewUserHandler(handlrServ, log)

		// Act
		err := handler.DeleteUser(ctx)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}
