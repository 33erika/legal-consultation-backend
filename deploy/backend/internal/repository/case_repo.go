package repository

import (
	"gorm.io/gorm"

	"legal-consultation/internal/models"
)

type CaseRepository struct {
	db *gorm.DB
}

func NewCaseRepository(db *gorm.DB) *CaseRepository {
	return &CaseRepository{db: db}
}

func (r *CaseRepository) Create(c *models.CaseCollection) error {
	return r.db.Create(c).Error
}

func (r *CaseRepository) GetByConsultationID(consultationID string) (*models.CaseCollection, error) {
	var c models.CaseCollection
	err := r.db.First(&c, "consultation_id = ?", consultationID).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CaseRepository) Delete(consultationID, collectorID string) error {
	return r.db.Where("consultation_id = ? AND collector_id = ?", consultationID, collectorID).
		Delete(&models.CaseCollection{}).Error
}

func (r *CaseRepository) ListByCollector(collectorID string, page, pageSize int) ([]models.CaseCollection, int64, error) {
	var cases []models.CaseCollection
	var total int64

	r.db.Model(&models.CaseCollection{}).Where("collector_id = ?", collectorID).Count(&total)
	err := r.db.
		Preload("Consultation").
		Preload("Consultation.Submitter").
		Where("collector_id = ?", collectorID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&cases).Error

	return cases, total, err
}

func (r *CaseRepository) Search(keyword, tag string, page, pageSize int) ([]models.CaseCollection, int64, error) {
	var cases []models.CaseCollection
	var total int64

	query := r.db.Model(&models.CaseCollection{})
	if keyword != "" {
		query = query.Joins("JOIN consultations ON consultations.id = case_collections.consultation_id").
			Where("consultations.title LIKE ? OR consultations.description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if tag != "" {
		query = query.Where("tags->'tags' ? ?", tag)
	}

	query.Count(&total)
	err := query.
		Preload("Consultation").
		Preload("Consultation.Submitter").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&cases).Error

	return cases, total, err
}
