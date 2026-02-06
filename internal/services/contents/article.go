package contents

import (
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/contents/entity"

	"gorm.io/gorm"
)

// ArticleService handles article operations
type ArticleService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Article]
}

// NewArticleService initializes the service
func NewArticleService(db *gorm.DB) *ArticleService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Article{})
	}
	return &ArticleService{
		db:   db,
		repo: model.NewCRUD[entity.Article](db),
	}
}

// Create adds a new article
func (s *ArticleService) Create(article *entity.Article) error {
	res := s.repo.Add(article)
	return res.Error
}

// Update modifies an existing article
func (s *ArticleService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Delete removes an article
func (s *ArticleService) Delete(id string) error {
	res := s.repo.Remove(id)
	return res.Error
}

// Get retrieves a single article by ID
func (s *ArticleService) Get(id string) (*entity.Article, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Article), nil
}

// List returns articles with filters
func (s *ArticleService) List(page, size int, filters map[string]interface{}) ([]entity.Article, int64, error) {
	// Handle special filters if necessary (e.g. searching by title)
	// For now, simple equality checks provided by CRUD are likely default, 
	// but if we need LIKE searches for Title, we might need custom logic or extended ListParams.
	// Assuming CRUD.List handles standard map filters. If we need "Title LIKE %query%", 
	// we'd add it to the query manually or use a helper. 
	// Given the context so far, we rely on the implementation of model.CRUD.
	// We'll pass standard filters.
	
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
	list := data["list"].([]entity.Article)
	total := data["total"].(int64)

	return list, total, nil
}

// IncrementView increases the view count for an article
func (s *ArticleService) IncrementView(id string) error {
	return s.db.Model(&entity.Article{}).Where("id = ?", id).UpdateColumn("view_times", gorm.Expr("view_times + ?", 1)).Error
}
