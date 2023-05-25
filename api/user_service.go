package api

import (
	"database/sql"
	"net/http"

	c "github.com/paquesqueue/bookstore/common"
	"github.com/paquesqueue/bookstore/utils"
)

type UserQueries interface {
	InsertUser(req RequestUser) (ResponseUser, error)
	SelectUser(username string) (ResponseUser, error)
	UpdateUser(username string, req RequestUser) (ResponseUser, error)
	DeleteUser(username string) error
}

type UserServices struct {
	query UserQueries
	log   c.Log
}

func NewUserService(q UserQueries, l c.Log) *UserServices {
	return &UserServices{q, l}
}

func (s UserServices) AddUser(req RequestUser) (ResponseUser, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Errorf("Error AddUser Hash Password : %v", err)
		return ResponseUser{}, &c.Err{Code: http.StatusInternalServerError, Remark: "Error AddUser Service", Original: err}
	}

	data := RequestUser{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Fullname: req.Fullname,
	}

	resp, err := s.query.InsertUser(data)
	if err != nil {
		s.log.Errorf("Error InsertUser : %v", err)
		return ResponseUser{}, &c.Err{Code: http.StatusInternalServerError, Remark: "Error AddUser Serivce", Original: err}
	}
	return resp, nil
}

func (s UserServices) GetUser(username string) (ResponseUser, error) {
	resp, err := s.query.SelectUser(username)
	if err != nil {
		s.log.Errorf("Error SelectUser : %v", err)
		switch err {
		case sql.ErrNoRows:
			return ResponseUser{}, &c.Err{Code: http.StatusNotFound, Remark: "Error User Not Found", Original: err}
		default:
			return ResponseUser{}, &c.Err{Code: http.StatusInternalServerError, Remark: "Error GetUser Service", Original: err}
		}
	}
	return resp, nil
}

func (s UserServices) PutUser(username string, req RequestUser) (ResponseUser, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Errorf("Error PutUser Hash Password : %v", err)
		return ResponseUser{}, &c.Err{Code: http.StatusInternalServerError, Remark: "Error PutUser Service", Original: err}
	}

	data := RequestUser{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Fullname: req.Fullname,
	}
	resp, err := s.query.UpdateUser(username, data)
	if err != nil {
		s.log.Errorf("Error UpdateUser : %v", err)
		return ResponseUser{}, &c.Err{Code: http.StatusInternalServerError, Remark: "Error PutUser Service", Original: err}
	}
	return resp, nil
}

func (s UserServices) DeleteUser(username string) error {
	err := s.query.DeleteUser(username)
	if err != nil {
		s.log.Errorf("Error DeleteUser : %v", err)
		return &c.Err{Code: http.StatusInternalServerError, Remark: "Error DeleteUser Service", Original: err}
	}
	return nil
}
