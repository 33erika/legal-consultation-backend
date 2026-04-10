package repository

import (
	"gorm.io/gorm"

	"legal-consultation/internal/models"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(attachment *models.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *AttachmentRepository) GetByID(id string) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.db.First(&attachment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *AttachmentRepository) Delete(id string) error {
	return r.db.Delete(&models.Attachment{}, "id = ?", id).Error
}

func (r *AttachmentRepository) ListByEntity(entityType, entityID string) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).Find(&attachments).Error
	return attachments, err
}

// 咨询附件关联
func (r *AttachmentRepository) CreateConsultationAttachment(ca *models.ConsultationAttachment) error {
	return r.db.Create(ca).Error
}

func (r *AttachmentRepository) ListConsultationAttachments(consultationID string) ([]models.ConsultationAttachment, error) {
	var cas []models.ConsultationAttachment
	err := r.db.
		Preload("Attachment").
		Where("consultation_id = ?", consultationID).
		Find(&cas).Error
	return cas, err
}

func (r *AttachmentRepository) DeleteConsultationAttachments(consultationID string) error {
	return r.db.Delete(&models.ConsultationAttachment{}, "consultation_id = ?", consultationID).Error
}

// 模板申请附件关联
func (r *AttachmentRepository) CreateTemplateRequestAttachment(ta *models.TemplateRequestAttachment) error {
	return r.db.Create(ta).Error
}

func (r *AttachmentRepository) ListTemplateRequestAttachments(requestID string) ([]models.TemplateRequestAttachment, error) {
	var tas []models.TemplateRequestAttachment
	err := r.db.
		Preload("Attachment").
		Where("template_request_id = ?", requestID).
		Find(&tas).Error
	return tas, err
}
