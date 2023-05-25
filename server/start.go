package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/paquesqueue/bookstore/common"
	"github.com/sirupsen/logrus"
)

func StartServer(s *http.Server, config common.Config) {
	log.Infof("Server starts on port %v", config.Port)
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down server %s", err)
		}
	}()
}

func GracefulShutdown(e *echo.Echo, log *logrus.Logger) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	<-signals

	log.Info("App is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error Server Shut Down : %s", err)
	}
	log.Info("Server is fully shutdown")
}
