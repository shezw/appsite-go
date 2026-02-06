package relation

import (
	"errors"

	"gorm.io/gorm"

	"appsite-go/internal/core/model"
	"appsite-go/internal/services/relation/entity"
)

// Service handles generic relation operations
type Service struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Relation]
}

// NewService initializes the service
func NewService(db *gorm.DB) *Service {
	if db != nil {
		_ = db.AutoMigrate(&entity.Relation{})
	}
	return &Service{
		db:   db,
		repo: model.NewCRUD[entity.Relation](db),
	}
}

// Bind creates a relationship
func (s *Service) Bind(item, itemType, target, targetType, relationType string) error {
	// Check exists
	var existing entity.Relation
	err := s.db.Where("item_id = ? AND item_type = ? AND relation_id = ? AND relation_type = ? AND type = ?",
		item, itemType, target, targetType, relationType).First(&existing).Error
	
	if err == nil {
		return nil // Already bound
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	rel := &entity.Relation{
		ItemID:       item,
		ItemType:     itemType,
		RelationID:   target,
		RelationType: targetType,
		Type:         relationType,
		Status:       "enabled",
	}
	return s.repo.Add(rel).Error
}

// Unbind removes a relationship
func (s *Service) Unbind(item, itemType, target, targetType, relationType string) error {
	return s.db.Where("item_id = ? AND item_type = ? AND relation_id = ? AND relation_type = ? AND type = ?",
		item, itemType, target, targetType, relationType).Delete(&entity.Relation{}).Error
}

// Check returns true if relationship exists
func (s *Service) Check(item, itemType, target, targetType, relationType string) bool {
	var count int64
	s.db.Model(&entity.Relation{}).Where("item_id = ? AND item_type = ? AND relation_id = ? AND relation_type = ? AND type = ?",
		item, itemType, target, targetType, relationType).Count(&count)
	return count > 0
}

// Count returns the number of relations (e.g., Follower count: target=User, type=follow)
func (s *Service) Count(target, targetType, relationType string) int64 {
	var count int64
	s.db.Model(&entity.Relation{}).Where("relation_id = ? AND relation_type = ? AND type = ?", target, targetType, relationType).Count(&count)
	return count
}

// ListRelations returns list of relations (e.g. following list)
func (s *Service) ListRelations(itemID, itemType, relationType string, page, size int) ([]entity.Relation, int64, error) {
	filters := map[string]interface{}{
		"item_id":   itemID,
		"item_type": itemType,
		"type":      relationType,
	}
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
	return data["list"].([]entity.Relation), data["total"].(int64), nil
}
