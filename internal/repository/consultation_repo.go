package repository

import (
	"time"

	"gorm.io/gorm"

	"legal-consultation/internal/models"
)

type ConsultationRepository struct {
	db *gorm.DB
}

func NewConsultationRepository(db *gorm.DB) *ConsultationRepository {
	return &ConsultationRepository{db: db}
}

func (r *ConsultationRepository) Create(consultation *models.Consultation) error {
	return r.db.Create(consultation).Error
}

func (r *ConsultationRepository) GetByID(id string) (*models.Consultation, error) {
	var consultation models.Consultation
	err := r.db.
		Preload("Submitter").
		Preload("Submitter.Department").
		Preload("Handler").
		First(&consultation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &consultation, nil
}

func (r *ConsultationRepository) Update(consultation *models.Consultation) error {
	return r.db.Save(consultation).Error
}

func (r *ConsultationRepository) Delete(id string) error {
	return r.db.Delete(&models.Consultation{}, "id = ?", id).Error
}

// 按状态列表
func (r *ConsultationRepository) ListByStatus(status string, page, pageSize int) ([]models.Consultation, int64, error) {
	var consultations []models.Consultation
	var total int64

	query := r.db.Model(&models.Consultation{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.
		Preload("Submitter").
		Preload("Handler").
		Order("CASE urgency WHEN 'very_urgent' THEN 1 WHEN 'urgent' THEN 2 ELSE 3 END").
		Order("submitted_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&consultations).Error

	return consultations, total, err
}

// 待处理咨询（咨询池）
func (r *ConsultationRepository) ListPending(page, pageSize int, urgency string) ([]models.Consultation, int64, error) {
	var consultations []models.Consultation
	var total int64

	query := r.db.Model(&models.Consultation{}).Where("status = ?", models.ConsultationStatusPending)
	if urgency != "" {
		query = query.Where("urgency = ?", urgency)
	}

	query.Count(&total)
	err := query.
		Preload("Submitter").
		Preload("Submitter.Department").
		Order("CASE urgency WHEN 'very_urgent' THEN 1 WHEN 'urgent' THEN 2 ELSE 3 END").
		Order("submitted_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&consultations).Error

	return consultations, total, err
}

// 我的待办（已接单）
func (r *ConsultationRepository) ListByHandler(handlerID string, status string, page, pageSize int) ([]models.Consultation, int64, error) {
	var consultations []models.Consultation
	var total int64

	query := r.db.Model(&models.Consultation{}).Where("handler_id = ?", handlerID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.
		Preload("Submitter").
		Order("accepted_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&consultations).Error

	return consultations, total, err
}

// 员工的咨询列表
func (r *ConsultationRepository) ListBySubmitter(submitterID string, status string, page, pageSize int) ([]models.Consultation, int64, error) {
	var consultations []models.Consultation
	var total int64

	query := r.db.Model(&models.Consultation{}).Where("submitter_id = ?", submitterID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.
		Preload("Handler").
		Order("submitted_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&consultations).Error

	return consultations, total, err
}

// 工作台统计
func (r *ConsultationRepository) GetStats() (map[string]interface{}, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))

	stats := make(map[string]interface{})

	// 今日新增
	var todayNew int64
	r.db.Model(&models.Consultation{}).Where("submitted_at >= ?", today).Count(&todayNew)
	stats["today_new"] = todayNew

	// 今日已回复
	var todayReplied int64
	r.db.Model(&models.Consultation{}).Where("first_replied_at >= ?", today).Count(&todayReplied)
	stats["today_replied"] = todayReplied

	// 待处理
	var pending int64
	r.db.Model(&models.Consultation{}).Where("status = ?", models.ConsultationStatusPending).Count(&pending)
	stats["pending"] = pending

	// 本周结案
	var weekClosed int64
	r.db.Model(&models.Consultation{}).Where("closed_at >= ?", weekStart).Count(&weekClosed)
	stats["week_closed"] = weekClosed

	return stats, nil
}

// 搜索
func (r *ConsultationRepository) Search(keyword string, category string, status string, page, pageSize int) ([]models.Consultation, int64, error) {
	var consultations []models.Consultation
	var total int64

	query := r.db.Model(&models.Consultation{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("consultation_type = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.
		Preload("Submitter").
		Preload("Handler").
		Order("submitted_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&consultations).Error

	return consultations, total, err
}

// 相似问题推荐
func (r *ConsultationRepository) FindSimilar(title string, limit int) ([]models.Consultation, error) {
	var consultations []models.Consultation
	err := r.db.
		Where("title LIKE ?", "%"+title[:min(len(title), 10)]+"%").
		Order("submitted_at DESC").
		Limit(limit).
		Find(&consultations).Error
	return consultations, err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
