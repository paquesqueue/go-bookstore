package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	c "github.com/paquesqueue/bookstore/common"
)

type UserHandlrQueries interface {
	AddUser(req RequestUser) (ResponseUser, error)
	GetUser(username string) (ResponseUser, error)
	PutUser(username string, req RequestUser) (ResponseUser, error)
	DeleteUser(username string) error
}

type UserHandlr struct {
	handler UserHandlrQueries
	log     c.Log
}

func NewUserHandler(h UserHandlrQueries, l c.Log) UserHandlr {
	return UserHandlr{h, l}
}

func (h UserHandlr) AddUser(ctx echo.Context) error {
	var req = RequestUser{}

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	resp, err := h.handler.AddUser(req)
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Error AddUser Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusCreated, resp)
}

func (h UserHandlr) GetUser(ctx echo.Context) error {
	username := ctx.Param("username")

	resp, err := h.handler.GetUser(username)
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Error GetUser Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (h UserHandlr) PutUser(ctx echo.Context) error {
	username := ctx.Param("username")
	var req = RequestUser{}
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	resp, err := h.handler.PutUser(username, req)
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Error PutUser Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (h UserHandlr) DeleteUser(ctx echo.Context) error {
	username := ctx.Param("username")

	err := h.handler.DeleteUser(username)
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Erro DelUser Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, "Deleted Successfully")
}
