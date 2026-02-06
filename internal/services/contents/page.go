package contents

import (
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/contents/entity"

	"gorm.io/gorm"
)

// PageService handles page operations
type PageService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Page]
}

// NewPageService initializes the service
func NewPageService(db *gorm.DB) *PageService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Page{})
	}
	return &PageService{
		db:   db,
		repo: model.NewCRUD[entity.Page](db),
	}
}

// Create adds a new page
func (s *PageService) Create(page *entity.Page) error {
	res := s.repo.Add(page)
	return res.Error
}

// Update modifies an existing page
func (s *PageService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Delete removes a page
func (s *PageService) Delete(id string) error {
	res := s.repo.Remove(id)
	return res.Error
}

// Get retrieves a single page by ID
func (s *PageService) Get(id string) (*entity.Page, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Page), nil
}

// GetByAlias retrieves a page by Alias
func (s *PageService) GetByAlias(alias string) (*entity.Page, error) {
	var page entity.Page
	if err := s.db.Where("alias = ?", alias).First(&page).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

// List returns pages with filters
func (s *PageService) List(page, size int, filters map[string]interface{}) ([]entity.Page, int64, error) {
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
	list := data["list"].([]entity.Page)
	total := data["total"].(int64)

	return list, total, nil
}

// IncrementView increases the view count for a page
func (s *PageService) IncrementView(id string) error {
	return s.db.Model(&entity.Page{}).Where("id = ?", id).UpdateColumn("view_times", gorm.Expr("view_times + ?", 1)).Error
}
