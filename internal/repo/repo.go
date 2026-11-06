package repo

import (
	"strings"
	"sync"
	"threats-info-exchange/internal/domain"
	"time"
)

type Repository struct {
	mu      sync.RWMutex
	entries map[string]domain.Threat
}

func NewRepository() *Repository {
	return &Repository{entries: make(map[string]domain.Threat)}
}

func (r *Repository) Get(md5 string) (*domain.Threat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.entries[strings.ToLower(md5)]

	if !ok {
		return nil, domain.ErrNotFound
	}

	return &item, nil
}

func (r *Repository) Create(req domain.CreateRequest) (*domain.Threat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := strings.ToLower(req.MD5)
	if _, ok := r.entries[key]; ok {
		return nil, domain.ErrExists
	}

	now := time.Now().UTC()
	item := domain.Threat{
		MD5:         key,
		Name:        req.Name,
		Description: strings.TrimSpace(req.Description),
		Category:    strings.TrimSpace(req.Category),
		Severity:    strings.TrimSpace(req.Severity),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.entries[key] = item
	return &item, nil
}

func (r *Repository) Update(md5 string, req domain.UpdateRequest) (*domain.Threat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := strings.ToLower(md5)
	item, ok := r.entries[key]
	if !ok {
		return nil, domain.ErrNotFound
	}

	item.Name = req.Name
	item.Description = strings.TrimSpace(req.Description)
	item.Category = strings.TrimSpace(req.Category)
	item.Severity = strings.TrimSpace(req.Severity)
	item.UpdatedAt = time.Now().UTC()
	r.entries[key] = item

	return &item, nil
}
