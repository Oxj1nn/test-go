package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("угроза не найдена")
	ErrExists   = errors.New("угроза уже существует")
)

type Threat struct {
	MD5         string    `json:"md5"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Severity    string    `json:"severity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	MD5         string `json:"md5"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
}

type UpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
}
