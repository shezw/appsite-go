package contents

import (
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/contents/entity"

	"gorm.io/gorm"
)

// TagService handles tag operations
type TagService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Tag]
}

// NewTagService initializes the service
func NewTagService(db *gorm.DB) *TagService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Tag{})
	}
	return &TagService{
		db:   db,
		repo: model.NewCRUD[entity.Tag](db),
	}
}

// Create adds a new tag
func (s *TagService) Create(tag *entity.Tag) error {
	res := s.repo.Add(tag)
	return res.Error
}

// Update modifies an existing tag
func (s *TagService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Delete removes a tag
func (s *TagService) Delete(id string) error {
	res := s.repo.Remove(id)
	return res.Error
}

// Get retrieves a single tag by ID
func (s *TagService) Get(id string) (*entity.Tag, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Tag), nil
}

// List returns tags with filters
func (s *TagService) List(page, size int, filters map[string]interface{}) ([]entity.Tag, int64, error) {
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
	list := data["list"].([]entity.Tag)
	total := data["total"].(int64)

	return list, total, nil
}
