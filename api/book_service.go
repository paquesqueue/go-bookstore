package api

import (
	"database/sql"
	"net/http"

	c "github.com/paquesqueue/bookstore/common"
)

type BookQueries interface {
	InsertBook(req RequestBook) (*ResponseBook, error)
	SelectAllBooks(params GetAllParams) ([]ResponseBook, error)
	SelectBookByID(id uint64) (*ResponseBook, error)
	UpdateBook(id uint64, req RequestBook) (*ResponseBook, error)
	DeleteBook(id uint64) error
}

type BookServices struct {
	query BookQueries
	log   c.Log
}

func NewBookService(s BookQueries, l c.Log) BookServices {
	return BookServices{s, l}
}

func (s BookServices) AddBook(req RequestBook) (*ResponseBook, error) {
	res, err := s.query.InsertBook(req)
	if err != nil {
		s.log.Errorf("Error InsertBook : %v", err)
		return nil, &c.Err{Code: http.StatusInternalServerError, Remark: "Error AddBook Service", Original: err}
	}
	return res, nil
}

func (s BookServices) ListAllBooks(params GetAllParams) ([]ResponseBook, error) {
	res, err := s.query.SelectAllBooks(params)
	if err != nil {
		s.log.Errorf("Error SelectAllBooks : %v", err)
		return nil, &c.Err{Code: http.StatusInternalServerError, Remark: "Error ListBooks Service", Original: err}
	}
	return res, nil
}

func (s BookServices) GetBookByID(id uint64) (*ResponseBook, error) {
	res, err := s.query.SelectBookByID(id)
	if err != nil {
		s.log.Errorf("Error SelectBookByID : %v", err)
		switch err {
		case sql.ErrNoRows:
			return nil, &c.Err{Code: http.StatusNotFound, Remark: "Error Book Not Found", Original: err}
		default:
			return nil, &c.Err{Code: http.StatusInternalServerError, Remark: "Error GetBook Service", Original: err}
		}
	}
	return res, nil
}

func (s BookServices) PutBook(id uint64, req RequestBook) (*ResponseBook, error) {
	res, err := s.query.UpdateBook(id, req)
	if err != nil {
		s.log.Errorf("Error UpdateBook : %v", err)
		return nil, &c.Err{Code: http.StatusInternalServerError, Remark: "Error UpdateBook Service", Original: err}
	}
	return res, nil
}

func (s BookServices) DelBook(id uint64) error {
	err := s.query.DeleteBook(id)
	if err != nil {
		s.log.Errorf("Error DeleteBook : %v", err)
		return &c.Err{Code: http.StatusInternalServerError, Remark: "Error DeleteBook Service", Original: err}
	}
	return nil
}
