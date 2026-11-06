package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusRecorder) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(body)
}

func (w *statusRecorder) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(recorder, r)

		if recorder.status == 0 {
			recorder.status = http.StatusOK
		}
		log.Info("HTTP-запрос", "метод", r.Method, "путь", r.URL.Path, "статус", recorder.status, "длительность", time.Since(start))
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Паника в HTTP-обработчике", "причина", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(map[string]string{"message": "Внутренняя ошибка сервера"}); err != nil {
					log.Error("Не удалось записать ответ об ошибке", "ошибка", err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
