package contents

import (
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/contents/entity"

	"gorm.io/gorm"
)

// BannerService handles banner operations
type BannerService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Banner]
}

// NewBannerService initializes the service
func NewBannerService(db *gorm.DB) *BannerService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Banner{})
	}
	return &BannerService{
		db:   db,
		repo: model.NewCRUD[entity.Banner](db),
	}
}

// Create adds a new banner
func (s *BannerService) Create(banner *entity.Banner) error {
	res := s.repo.Add(banner)
	return res.Error
}

// Update modifies an existing banner
func (s *BannerService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Delete removes a banner
func (s *BannerService) Delete(id string) error {
	res := s.repo.Remove(id)
	return res.Error
}

// Get retrieves a single banner by ID
func (s *BannerService) Get(id string) (*entity.Banner, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Banner), nil
}

// List returns banners with filters
func (s *BannerService) List(page, size int, filters map[string]interface{}) ([]entity.Banner, int64, error) {
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
	list := data["list"].([]entity.Banner)
	total := data["total"].(int64)

	return list, total, nil
}

// IncrementClick increases the click count for a banner
func (s *BannerService) IncrementClick(id string) error {
	return s.db.Model(&entity.Banner{}).Where("id = ?", id).UpdateColumn("click_times", gorm.Expr("click_times + ?", 1)).Error
}

// IncrementView increases the view count for a banner
func (s *BannerService) IncrementView(id string) error {
	return s.db.Model(&entity.Banner{}).Where("id = ?", id).UpdateColumn("view_times", gorm.Expr("view_times + ?", 1)).Error
}
