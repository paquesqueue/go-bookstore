package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/paquesqueue/bookstore/api"
	"github.com/paquesqueue/bookstore/common"
	"github.com/paquesqueue/bookstore/server"
)

func main() {
	log := common.InitLog()

	config := common.InitConfig()
	log.Info("Success Config Loaded")

	db, err := api.InitDB(config, log)
	if err != nil {
		log.Fatal("Error Database Init Failed")
	}
	log.Info("Success Database Connection Initialized")

	echo := echo.New()

	reqLog := common.InitRequestLog()
	
	server.InitMiddleware(echo, reqLog, config)
	server.InitRoutes(echo, db, log)

	serv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: echo,
	}

	server.StartServer(serv, config)
	server.GracefulShutdown(echo, log)
}
