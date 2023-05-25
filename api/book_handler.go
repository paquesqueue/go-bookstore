package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	c "github.com/paquesqueue/bookstore/common"
)

type BookHandlrQueries interface {
	AddBook(req RequestBook) (*ResponseBook, error)
	ListAllBooks(params GetAllParams) ([]ResponseBook, error)
	GetBookByID(id uint64) (*ResponseBook, error)
	PutBook(id uint64, req RequestBook) (*ResponseBook, error)
	DelBook(id uint64) error
}

type BookHandlr struct {
	handler BookHandlrQueries
	log     c.Log
}

func NewBookHandlr(h BookHandlrQueries, l c.Log) BookHandlr {
	return BookHandlr{h, l}
}

func (h BookHandlr) AddBook(ctx echo.Context) error {
	req := RequestBook{}
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	res, err := h.handler.AddBook(req)
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Error AddBook Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusCreated, res)
}

func (h BookHandlr) ListAllBooks(ctx echo.Context) error {
	var req RequestGetAll
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	params := GetAllParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	res, err := h.handler.ListAllBooks(params)
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Error ListAllBooks Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, res)
}

func (h BookHandlr) GetBookByID(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	res, err := h.handler.GetBookByID(uint64(id))
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Error GetBookByID Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, res)
}

func (h BookHandlr) PutBook(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	req := RequestBook{}
	err = ctx.Bind(&req)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	res, err := h.handler.PutBook(uint64(id), req)
	if err != nil {
		if cmErrm, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErrm.Code)
		}
		h.log.Errorf("Error PutBook Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, res)
}

func (h BookHandlr) DelBook(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	err = h.handler.DelBook(uint64(id))
	if err != nil {
		if cmErr, ok := err.(*c.Err); ok {
			return ctx.NoContent(cmErr.Code)
		}
		h.log.Errorf("Error DelBook Handler : %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, "Deleted Successfully")
}
