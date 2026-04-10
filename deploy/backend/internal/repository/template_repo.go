package repository

import (
	"time"

	"gorm.io/gorm"

	"legal-consultation/internal/models"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

// ============ 模板申请 ============

func (r *TemplateRepository) CreateRequest(req *models.TemplateRequest) error {
	return r.db.Create(req).Error
}

func (r *TemplateRepository) GetRequestByID(id string) (*models.TemplateRequest, error) {
	var req models.TemplateRequest
	err := r.db.
		Preload("Submitter").
		Preload("L1Approver").
		Preload("Drafter").
		Preload("Reviewer").
		Preload("ExistingTemplate").
		Preload("Attachments").
		Preload("Attachments.Attachment").
		Preload("ApprovalLogs").
		Preload("ApprovalLogs.Approver").
		First(&req, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *TemplateRepository) UpdateRequest(req *models.TemplateRequest) error {
	return r.db.Save(req).Error
}

func (r *TemplateRepository) DeleteRequest(id string) error {
	return r.db.Delete(&models.TemplateRequest{}, "id = ?", id).Error
}

func (r *TemplateRepository) ListRequestsBySubmitter(submitterID string, page, pageSize int) ([]models.TemplateRequest, int64, error) {
	var requests []models.TemplateRequest
	var total int64

	r.db.Model(&models.TemplateRequest{}).Where("submitter_id = ?", submitterID).Count(&total)
	err := r.db.
		Where("submitter_id = ?", submitterID).
		Order("submitted_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

func (r *TemplateRepository) ListRequestsByStatus(status string, page, pageSize int) ([]models.TemplateRequest, int64, error) {
	var requests []models.TemplateRequest
	var total int64

	query := r.db.Model(&models.TemplateRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.
		Preload("Submitter").
		Order("submitted_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

// 待审批（业务主管视角）
func (r *TemplateRepository) ListPendingApproval(page, pageSize int) ([]models.TemplateRequest, int64, error) {
	var requests []models.TemplateRequest
	var total int64

	r.db.Model(&models.TemplateRequest{}).
		Where("status = ?", models.TemplateRequestStatusPendingApproval).
		Count(&total)

	err := r.db.
		Where("status = ?", models.TemplateRequestStatusPendingApproval).
		Preload("Submitter").
		Order("submitted_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

// 待拟写（法务专员视角）
func (r *TemplateRepository) ListPendingDraft(page, pageSize int) ([]models.TemplateRequest, int64, error) {
	var requests []models.TemplateRequest
	var total int64

	r.db.Model(&models.TemplateRequest{}).
		Where("status = ?", models.TemplateRequestStatusPendingDraft).
		Count(&total)

	err := r.db.
		Where("status = ?", models.TemplateRequestStatusPendingDraft).
		Preload("Submitter").
		Order("submitted_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

// 待审核（法务负责人视角）
func (r *TemplateRepository) ListPendingReview(page, pageSize int) ([]models.TemplateRequest, int64, error) {
	var requests []models.TemplateRequest
	var total int64

	r.db.Model(&models.TemplateRequest{}).
		Where("status = ?", models.TemplateRequestStatusPendingReview).
		Count(&total)

	err := r.db.
		Where("status = ?", models.TemplateRequestStatusPendingReview).
		Preload("Submitter").
		Preload("Drafter").
		Order("drafted_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&requests).Error

	return requests, total, err
}

// 模板申请统计
func (r *TemplateRepository) GetRequestStats() (map[string]interface{}, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	stats := make(map[string]interface{})

	var pendingDraft int64
	r.db.Model(&models.TemplateRequest{}).Where("status = ?", models.TemplateRequestStatusPendingDraft).Count(&pendingDraft)
	stats["pending_draft"] = pendingDraft

	var pendingReview int64
	r.db.Model(&models.TemplateRequest{}).Where("status = ?", models.TemplateRequestStatusPendingReview).Count(&pendingReview)
	stats["pending_review"] = pendingReview

	var todayNew int64
	r.db.Model(&models.TemplateRequest{}).Where("submitted_at >= ?", today).Count(&todayNew)
	stats["today_new_requests"] = todayNew

	return stats, nil
}

// 创建审批日志
func (r *TemplateRepository) CreateApprovalLog(log *models.TemplateApprovalLog) error {
	return r.db.Create(log).Error
}

// ============ 合同模板 ============

func (r *TemplateRepository) CreateTemplate(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *TemplateRepository) GetTemplateByID(id string) (*models.Template, error) {
	var template models.Template
	err := r.db.
		Preload("PublishedBy").
		First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepository) GetTemplateByName(name string) (*models.Template, error) {
	var template models.Template
	// 获取最新版本
	err := r.db.
		Where("name = ? AND status = ?", name, models.TemplateStatusPublished).
		Order("created_at DESC").
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepository) UpdateTemplate(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *TemplateRepository) ListTemplates(contractType string, keyword string, page, pageSize int) ([]models.Template, int64, error) {
	var templates []models.Template
	var total int64

	query := r.db.Model(&models.Template{}).Where("status = ?", models.TemplateStatusPublished)
	if contractType != "" {
		query = query.Where("contract_type = ?", contractType)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.
		Preload("PublishedBy").
		Order("published_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&templates).Error

	return templates, total, err
}

// 获取模板版本历史
func (r *TemplateRepository) GetTemplateVersions(templateID string) ([]models.TemplateVersion, error) {
	var versions []models.TemplateVersion
	err := r.db.
		Where("template_id = ?", templateID).
		Order("published_at DESC").
		Find(&versions).Error
	return versions, err
}

func (r *TemplateRepository) CreateTemplateVersion(version *models.TemplateVersion) error {
	return r.db.Create(version).Error
}

// 版本对比
func (r *TemplateRepository) GetVersionByID(versionID string) (*models.TemplateVersion, error) {
	var version models.TemplateVersion
	err := r.db.First(&version, "id = ?", versionID).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// 增加下载次数
func (r *TemplateRepository) IncrementDownloadCount(id string) error {
	return r.db.Model(&models.Template{}).Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + ?", 1)).Error
}

// 禁用/启用模板
func (r *TemplateRepository) ToggleTemplateStatus(id string, status string) error {
	return r.db.Model(&models.Template{}).Where("id = ?", id).Update("status", status).Error
}

// 模板统计
func (r *TemplateRepository) GetTemplateStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var total int64
	r.db.Model(&models.Template{}).Where("status = ?", models.TemplateStatusPublished).Count(&total)
	stats["total_published"] = total

	return stats, nil
}
