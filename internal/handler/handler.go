package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"

	"threats-info-exchange/internal/domain"
)

var md5Pattern = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

const (
	messageInvalidGetMD5    = "Некорректный формат MD5"
	messageInvalidUpdate    = "Некорректные данные запроса или MD5"
	messageInvalidCreate    = "Некорректные данные запроса"
	messageThreatNotFound   = "Угроза с указанным MD5 не найден"
	messageInternalError    = "Внутренняя ошибка сервера"
	messageThreatExists     = "Угроза с таким MD5 уже существует"
	messageMethodNotAllowed = "Метод не поддерживается"
	messageRouteNotFound    = "Ресурс не найден"

	detailsInvalidMD5      = "MD5 должен содержать ровно 32 шестнадцатеричных символа"
	detailsInvalidSeverity = "Поле severity должно иметь одно из значений: unknown, low, medium, high, critical"
)

type Service interface {
	Get(md5 string) (*domain.Threat, error)
	Create(req domain.CreateRequest) (*domain.Threat, error)
	Update(md5 string, req domain.UpdateRequest) (*domain.Threat, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func writeError(w http.ResponseWriter, status int, msg, details string) {
	if status >= http.StatusInternalServerError {
		log.Error("Ошибка HTTP", "статус", status, "сообщение", msg, "подробности", details)
	} else {
		log.Warn("Ошибка HTTP", "статус", status, "сообщение", msg, "подробности", details)
	}

	type errorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse{
		Code:    status,
		Message: msg,
		Details: details,
	}); err != nil {
		log.Error("Не удалось записать ответ об ошибке", "ошибка", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("Не удалось записать ответ JSON", "статус", status, "ошибка", err)
	}
}

// GET /info/{md5}
func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, messageMethodNotAllowed, "")
		return
	}

	md5 := strings.TrimPrefix(r.URL.Path, "/info/")
	if md5 == "" || strings.Contains(md5, "/") {
		writeError(w, http.StatusNotFound, messageRouteNotFound, "")
		return
	}

	if !md5Pattern.MatchString(md5) {
		writeError(w, http.StatusBadRequest, messageInvalidGetMD5, detailsInvalidMD5)
		return
	}

	item, err := h.svc.Get(md5)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, messageThreatNotFound, fmt.Sprintf("Угроза с MD5 %s не найдена", md5))
		} else {
			writeError(w, http.StatusInternalServerError, messageInternalError, "")
		}
		return
	}

	writeJSON(w, http.StatusOK, item)
}

// POST /info/
func (h *Handler) PostInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, messageMethodNotAllowed, "")
		return
	}

	var req domain.CreateRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, messageInvalidCreate, err.Error())
		return
	}

	if !md5Pattern.MatchString(req.MD5) {
		writeError(w, http.StatusBadRequest, messageInvalidCreate, detailsInvalidMD5)
		return
	}

	if !validSeverity(req.Severity) {
		writeError(w, http.StatusBadRequest, messageInvalidCreate, detailsInvalidSeverity)
		return
	}

	item, err := h.svc.Create(req)
	if err != nil {
		if errors.Is(err, domain.ErrExists) {
			writeError(w, http.StatusConflict, messageThreatExists, fmt.Sprintf("Угроза с MD5 %s уже существует", req.MD5))
		} else {
			writeError(w, http.StatusInternalServerError, messageInternalError, "")
		}
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

// PUT /info/{md5}
func (h *Handler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, messageMethodNotAllowed, "")
		return
	}

	md5 := strings.TrimPrefix(r.URL.Path, "/info/")
	if md5 == "" || strings.Contains(md5, "/") {
		writeError(w, http.StatusNotFound, messageRouteNotFound, "")
		return
	}

	if !md5Pattern.MatchString(md5) {
		writeError(w, http.StatusBadRequest, messageInvalidUpdate, detailsInvalidMD5)
		return
	}

	var req domain.UpdateRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, messageInvalidUpdate, err.Error())
		return
	}

	if !validSeverity(req.Severity) {
		writeError(w, http.StatusBadRequest, messageInvalidUpdate, detailsInvalidSeverity)
		return
	}

	item, err := h.svc.Update(md5, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, messageThreatNotFound, fmt.Sprintf("Угроза с MD5 %s не найдена", md5))
		} else {
			writeError(w, http.StatusInternalServerError, messageInternalError, "")
		}
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func validSeverity(value string) bool {
	switch value {
	case "", "unknown", "low", "medium", "high", "critical":
		return true
	default:
		return false
	}
}
