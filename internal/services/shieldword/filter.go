package shieldword

import (
	"log"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"appsite-go/internal/core/model"
	"appsite-go/internal/services/shieldword/entity"
)

// Service handles shield word operations and filtering
type Service struct {
	db        *gorm.DB
	repo      *model.CRUD[entity.Word]
	cache     []string
	cacheTime time.Time
	mu        sync.RWMutex
}

// NewService initializes the service
func NewService(db *gorm.DB) *Service {
	if db != nil {
		_ = db.AutoMigrate(&entity.Word{})
	}
	return &Service{
		db:   db,
		repo: model.NewCRUD[entity.Word](db),
	}
}

// ensureCache loads words into memory if expired
func (s *Service) ensureCache() {
	s.mu.RLock()
	// Cache for 5 minutes
	if time.Since(s.cacheTime) < 5*time.Minute && s.cache != nil {
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check
	if time.Since(s.cacheTime) < 5*time.Minute && s.cache != nil {
		return
	}

	var words []entity.Word
	if err := s.db.Where("status = ?", "enabled").Find(&words).Error; err != nil {
		log.Printf("Failed to load shield words: %v", err)
		return
	}

	s.cache = make([]string, len(words))
	for i, w := range words {
		s.cache[i] = w.Title
	}
	s.cacheTime = time.Now()
}

// Create adds a new shield word
func (s *Service) Create(word *entity.Word) error {
	res := s.repo.Add(word)
	if res.Success {
		s.mu.Lock()
		s.cache = nil // Invalidate cache
		s.mu.Unlock()
	}
	return res.Error
}

// Delete removes a shield word
func (s *Service) Delete(id string) error {
	res := s.repo.Remove(id)
	if res.Success {
		s.mu.Lock()
		s.cache = nil // Invalidate cache
		s.mu.Unlock()
	}
	return res.Error
}

// Check returns true if content contains sensitive words
func (s *Service) Check(content string) bool {
	s.ensureCache()
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, word := range s.cache {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}

// Replace masks sensitive words in content
func (s *Service) Replace(content string) string {
	s.ensureCache()
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, word := range s.cache {
		if word == "" {
			continue
		}
		mask := strings.Repeat("*", len([]rune(word)))
		content = strings.ReplaceAll(content, word, mask)
	}
	return content
}
