package contents

import (
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/contents/entity"

	"gorm.io/gorm"
)

// CommentService handles comment operations
type CommentService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Comment]
}

// NewCommentService initializes the service
func NewCommentService(db *gorm.DB) *CommentService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Comment{})
	}
	return &CommentService{
		db:   db,
		repo: model.NewCRUD[entity.Comment](db),
	}
}

// Create adds a new comment
func (s *CommentService) Create(comment *entity.Comment) error {
	res := s.repo.Add(comment)
	return res.Error
}

// Update modifies an existing comment
func (s *CommentService) Update(id string, updates map[string]interface{}) error {
	res := s.repo.Update(id, updates)
	return res.Error
}

// Delete removes a comment
func (s *CommentService) Delete(id string) error {
	res := s.repo.Remove(id)
	return res.Error
}

// Get retrieves a single comment by ID
func (s *CommentService) Get(id string) (*entity.Comment, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Comment), nil
}

// List returns comments with filters
func (s *CommentService) List(page, size int, filters map[string]interface{}) ([]entity.Comment, int64, error) {
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
	list := data["list"].([]entity.Comment)
	total := data["total"].(int64)

	return list, total, nil
}
