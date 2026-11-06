package main

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"

	"threats-info-exchange/internal/handler"
	"threats-info-exchange/internal/middleware"
	"threats-info-exchange/internal/repo"
	"threats-info-exchange/internal/router"
)

func main() {
	h := handler.NewHandler(repo.NewRepository())
	app := middleware.Logging(middleware.Recover(router.New(h)))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          log.Default().StandardLog(log.StandardLogOptions{ForceLevel: log.ErrorLevel}),
	}

	log.Info("Сервер слушает соединения", "адрес", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Сервер остановился с ошибкой", "ошибка", err)
	}
}
