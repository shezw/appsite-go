package contents

import (
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/contents/entity"

	"gorm.io/gorm"
)

// MediaService handles media operations
type MediaService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Media]
}

// NewMediaService initializes the service
func NewMediaService(db *gorm.DB) *MediaService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Media{})
	}
	return &MediaService{
		db:   db,
		repo: model.NewCRUD[entity.Media](db),
	}
}

// Create adds a new media record
func (s *MediaService) Create(media *entity.Media) error {
	res := s.repo.Add(media)
	return res.Error
}

// Update modifies an existing media record
func (s *MediaService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Delete removes a media record
func (s *MediaService) Delete(id string) error {
	res := s.repo.Remove(id)
	return res.Error
}

// Get retrieves a single media record by ID
func (s *MediaService) Get(id string) (*entity.Media, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Media), nil
}

// List returns media records with filters
func (s *MediaService) List(page, size int, filters map[string]interface{}) ([]entity.Media, int64, error) {
	res := s.repo.List(&model.ListParams{
		Page:     page,
		PageSize: size,
		Filters:  filters,
		Sort:     "sort desc, created_at desc",
	})

	if !res.Success {
		return nil, 0, res.Error
	}

	data := res.Data.(map[string]interface{})
	list := data["list"].([]entity.Media)
	total := data["total"].(int64)

	return list, total, nil
}
