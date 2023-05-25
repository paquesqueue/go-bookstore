package server

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	api "github.com/paquesqueue/bookstore/api"
	"github.com/sirupsen/logrus"
)

func InitRoutes(e *echo.Echo, dbConn *sql.DB, log *logrus.Logger) {
	conn := api.NewDB(dbConn)

	bookServ := api.NewBookService(conn, log)
	bookHandlr := api.NewBookHandlr(bookServ, log)

	e.POST("/books", bookHandlr.AddBook)
	e.GET("/books", bookHandlr.ListAllBooks)
	e.GET("/books/:id", bookHandlr.GetBookByID)
	e.PUT("/books/:id", bookHandlr.PutBook)
	e.DELETE("/books/:id", bookHandlr.DelBook)

	userServ := api.NewUserService(conn, log)
	userHandlr := api.NewUserHandler(userServ, log)

	e.POST("/users", userHandlr.AddUser)
	e.GET("/users/:username", userHandlr.GetUser)
	e.PUT("/users/:username", userHandlr.PutUser)
	e.DELETE("/users/:username", userHandlr.DeleteUser)
}
