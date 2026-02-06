package form

import (
	"gorm.io/gorm"

	"appsite-go/internal/core/model"
	"appsite-go/internal/services/form/entity"
)

// SubmissionService handles form submissions
type SubmissionService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Request]
}

// NewSubmissionService initializes the service
func NewSubmissionService(db *gorm.DB) *SubmissionService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Request{})
	}
	return &SubmissionService{
		db:   db,
		repo: model.NewCRUD[entity.Request](db),
	}
}

// Submit creates a new form request
func (s *SubmissionService) Submit(req *entity.Request) error {
	req.Status = "pending"
	res := s.repo.Add(req)
	return res.Error
}

// Review updates the status of a request
func (s *SubmissionService) Review(id string, approved bool) error {
	status := "rejected"
	if approved {
		status = "applied"
	}

	res := s.repo.Update(id, map[string]interface{}{
		"status": status,
	})
	
	if res.Success {
		// In a real system, we would trigger the ApplyCall/RejectCall hooks here
		// For now, just recording the state change is sufficient for Phase 3
	}

	return res.Error
}

// Get retrieves a single request
func (s *SubmissionService) Get(id string) (*entity.Request, error) {
	res := s.repo.Get(id)
	if !res.Success {
		return nil, res.Error
	}
	return res.Data.(*entity.Request), nil
}

// List retrieves requests with filters
func (s *SubmissionService) List(page, size int, filters map[string]interface{}) ([]entity.Request, int64, error) {
	res := s.repo.List(&model.ListParams{
		Page:     page,
		PageSize: size,
		Filters:  filters,
		Sort:     "created_at desc",
	})

	if !res.Success {
		return nil, 0, res.Error
	}

	data := res.Data.(map[string]interface{})
	list := data["list"].([]entity.Request)
	total := data["total"].(int64)

	return list, total, nil
}
