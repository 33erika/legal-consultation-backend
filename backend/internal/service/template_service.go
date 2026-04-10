package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"legal-consultation/internal/models"
	"legal-consultation/internal/repository"
	"legal-consultation/internal/utils"
)

type TemplateService struct {
	templateRepo     *repository.TemplateRepository
	attachmentRepo   *repository.AttachmentRepository
	notificationSvc  *NotificationService
}

func NewTemplateService(
	templateRepo *repository.TemplateRepository,
	attachmentRepo *repository.AttachmentRepository,
	notificationSvc *NotificationService,
) *TemplateService {
	return &TemplateService{
		templateRepo:    templateRepo,
		attachmentRepo:  attachmentRepo,
		notificationSvc: notificationSvc,
	}
}

// ============ 模板申请 ============

// 创建申请
func (s *TemplateService) CreateTemplateRequest(req *CreateTemplateRequestRequest, submitterID string) (*models.TemplateRequest, error) {
	templateReq := &models.TemplateRequest{
		ID:               uuid.New().String(),
		RequestNo:        utils.GenerateTemplateRequestNo(),
		RequestType:      req.RequestType,
		ExistingTemplateID: req.ExistingTemplateID,
		ContractType:     req.ContractType,
		Title:            req.Title,
		BusinessScenario: req.BusinessScenario,
		BusinessFlow:     req.BusinessFlow,
		KeyClauses:       req.KeyClauses,
		DiffFromExisting: req.DiffFromExisting,
		ReferenceFiles:   req.ReferenceFiles,
		ExpectedDate:     req.ExpectedDate,
		Status:           models.TemplateRequestStatusPendingApproval,
		CurrentStep:      1,
		SubmitterID:      submitterID,
		SubmittedAt:      time.Now(),
	}

	if err := s.templateRepo.CreateRequest(templateReq); err != nil {
		return nil, err
	}

	// 通知业务主管
	s.notificationSvc.NotifyTemplateRequestPending(templateReq)

	return templateReq, nil
}

type CreateTemplateRequestRequest struct {
	RequestType        string                  `json:"request_type" binding:"required,oneof=new update"`
	ExistingTemplateID *string                  `json:"existing_template_id"`
	ContractType       string                   `json:"contract_type" binding:"required"`
	Title              string                   `json:"title" binding:"required,max=100"`
	BusinessScenario   string                   `json:"business_scenario"`
	BusinessFlow       string                   `json:"business_flow"`
	KeyClauses         string                   `json:"key_clauses"`
	DiffFromExisting   string                   `json:"diff_from_existing"`
	ReferenceFiles     models.JSONType          `json:"reference_files"`
	ExpectedDate       *time.Time               `json:"expected_date"`
}

// 获取申请详情
func (s *TemplateService) GetTemplateRequest(id string) (*models.TemplateRequest, error) {
	return s.templateRepo.GetRequestByID(id)
}

// 列表 - 员工视角
func (s *TemplateService) ListMyRequests(submitterID string, page, pageSize int) ([]models.TemplateRequest, int64, error) {
	return s.templateRepo.ListRequestsBySubmitter(submitterID, page, pageSize)
}

// 列表 - 待审批
func (s *TemplateService) ListPendingApproval(page, pageSize int) ([]models.TemplateRequest, int64, error) {
	return s.templateRepo.ListPendingApproval(page, pageSize)
}

// 列表 - 待拟写
func (s *TemplateService) ListPendingDraft(page, pageSize int) ([]models.TemplateRequest, int64, error) {
	return s.templateRepo.ListPendingDraft(page, pageSize)
}

// 列表 - 待审核
func (s *TemplateService) ListPendingReview(page, pageSize int) ([]models.TemplateRequest, int64, error) {
	return s.templateRepo.ListPendingReview(page, pageSize)
}

// L1 审批
func (s *TemplateService) ApproveRequest(id string, req *ApproveTemplateRequest, approverID string) error {
	templateReq, err := s.templateRepo.GetRequestByID(id)
	if err != nil {
		return errors.New("申请不存在")
	}

	if templateReq.Status != models.TemplateRequestStatusPendingApproval {
		return errors.New("该申请不在待审批状态")
	}

	switch req.Action {
	case "approve":
		templateReq.Status = models.TemplateRequestStatusPendingDraft
		templateReq.CurrentStep = 2
		templateReq.L1ApproverID = &approverID
		now := time.Now()
		templateReq.L1ApprovedAt = &now
		templateReq.L1Comment = req.Comment

		// 通知法务专员
		s.notificationSvc.NotifyTemplateRequestPendingDraft(templateReq)

	case "reject":
		templateReq.Status = models.TemplateRequestStatusRejected
		templateReq.L1ApproverID = &approverID
		now := time.Now()
		templateReq.L1ApprovedAt = &now
		templateReq.L1Comment = req.Comment

		// 通知申请人
		s.notificationSvc.NotifyTemplateRequestRejected(templateReq)

	case "return_for_supplement":
		templateReq.Status = models.TemplateRequestStatusSupplemented
		templateReq.L1ApproverID = &approverID
		now := time.Now()
		templateReq.L1ApprovedAt = &now
		templateReq.L1Comment = req.Comment

		// 通知申请人
		s.notificationSvc.NotifyTemplateRequestReturnForSupplement(templateReq)
	}

	// 记录审批日志
	approvalLog := &models.TemplateApprovalLog{
		ID:               uuid.New().String(),
		TemplateRequestID: id,
		ApproverID:       approverID,
		Action:           req.Action,
		Comment:          req.Comment,
		CreatedAt:        time.Now(),
	}
	s.templateRepo.CreateApprovalLog(approvalLog)

	return s.templateRepo.UpdateRequest(templateReq)
}

type ApproveTemplateRequest struct {
	Action  string `json:"action" binding:"required,oneof=approve reject return_for_supplement"`
	Comment string `json:"comment"`
}

// 拟写模板
func (s *TemplateService) DraftTemplate(id string, req *DraftTemplateRequest, drafterID string) error {
	templateReq, err := s.templateRepo.GetRequestByID(id)
	if err != nil {
		return errors.New("申请不存在")
	}

	if templateReq.Status != models.TemplateRequestStatusPendingDraft {
		return errors.New("该申请不在待拟写状态")
	}

	// 创建模板
	template := &models.Template{
		ID:               uuid.New().String(),
		Name:             req.Name,
		ContractType:     templateReq.ContractType,
		Version:          "v1.0",
		Description:      req.Description,
		FilePath:         req.FilePath,
		EditableClauses:  req.EditableClauses,
		Status:           models.TemplateStatusDraft,
		TemplateRequestID: &id,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.templateRepo.CreateTemplate(template); err != nil {
		return err
	}

	// 更新申请状态
	templateReq.Status = models.TemplateRequestStatusPendingReview
	templateReq.CurrentStep = 3
	templateReq.DrafterID = &drafterID
	now := time.Now()
	templateReq.DraftedAt = &now

	if err := s.templateRepo.UpdateRequest(templateReq); err != nil {
		return err
	}

	// 通知法务负责人
	s.notificationSvc.NotifyTemplateRequestPendingReview(templateReq)

	return nil
}

type DraftTemplateRequest struct {
	Name           string `json:"name" binding:"required,max=100"`
	Description    string `json:"description" binding:"required,max=2000"`
	FilePath       string `json:"file_path" binding:"required"`
	EditableClauses string `json:"editable_clauses"`
}

// 保存草稿
func (s *TemplateService) SaveDraft(id string, drafterID string) error {
	templateReq, err := s.templateRepo.GetRequestByID(id)
	if err != nil {
		return errors.New("申请不存在")
	}

	// 草稿状态不变，只更新时间
	templateReq.DrafterID = &drafterID
	return s.templateRepo.UpdateRequest(templateReq)
}

// 审核
func (s *TemplateService) ReviewTemplate(id string, req *ReviewTemplateRequest, reviewerID string) error {
	templateReq, err := s.templateRepo.GetRequestByID(id)
	if err != nil {
		return errors.New("申请不存在")
	}

	if templateReq.Status != models.TemplateRequestStatusPendingReview {
		return errors.New("该申请不在待审核状态")
	}

	switch req.Action {
	case "approve":
		templateReq.Status = models.TemplateRequestStatusPublished
		templateReq.ReviewerID = &reviewerID
		now := time.Now()
		templateReq.ReviewedAt = &now
		templateReq.ReviewComment = req.Comment

		// 更新模板状态为已发布
		template, err := s.templateRepo.GetTemplateByID(req.TemplateID)
		if err == nil {
			template.Status = models.TemplateStatusPublished
			template.PublishedByID = &reviewerID
			template.PublishedAt = &now
			s.templateRepo.UpdateTemplate(template)
		}

		// 通知全员
		s.notificationSvc.NotifyTemplatePublished(templateReq)

	case "return_for_modification":
		templateReq.Status = models.TemplateRequestStatusPendingDraft
		templateReq.CurrentStep = 2
		templateReq.ReviewerID = &reviewerID
		now := time.Now()
		templateReq.ReviewedAt = &now
		templateReq.ReviewComment = req.Comment

		// 通知法务专员
		s.notificationSvc.NotifyTemplateRequestReturnForModification(templateReq)
	}

	// 记录审批日志
	approvalLog := &models.TemplateApprovalLog{
		ID:               uuid.New().String(),
		TemplateRequestID: id,
		ApproverID:       reviewerID,
		Action:           req.Action,
		Comment:          req.Comment,
		CreatedAt:        time.Now(),
	}
	s.templateRepo.CreateApprovalLog(approvalLog)

	return s.templateRepo.UpdateRequest(templateReq)
}

type ReviewTemplateRequest struct {
	Action     string `json:"action" binding:"required,oneof=approve return_for_modification"`
	Comment    string `json:"comment"`
	TemplateID string `json:"template_id" binding:"required"`
}

// 统计
func (s *TemplateService) GetRequestStats() (map[string]interface{}, error) {
	return s.templateRepo.GetRequestStats()
}

// ============ 合同模板库 ============

// 列表
func (s *TemplateService) ListTemplates(contractType, keyword string, page, pageSize int) ([]models.Template, int64, error) {
	return s.templateRepo.ListTemplates(contractType, keyword, page, pageSize)
}

// 详情
func (s *TemplateService) GetTemplate(id string) (*models.Template, error) {
	return s.templateRepo.GetTemplateByID(id)
}

// 下载
func (s *TemplateService) DownloadTemplate(id string) error {
	return s.templateRepo.IncrementDownloadCount(id)
}

// 版本历史
func (s *TemplateService) GetTemplateVersions(templateID string) ([]models.TemplateVersion, error) {
	return s.templateRepo.GetTemplateVersions(templateID)
}

// 版本对比
func (s *TemplateService) CompareVersions(versionAID, versionBID string) (*models.TemplateVersion, *models.TemplateVersion, error) {
	versionA, err := s.templateRepo.GetVersionByID(versionAID)
	if err != nil {
		return nil, nil, errors.New("版本A不存在")
	}

	versionB, err := s.templateRepo.GetVersionByID(versionBID)
	if err != nil {
		return nil, nil, errors.New("版本B不存在")
	}

	return versionA, versionB, nil
}

// 禁用/启用
func (s *TemplateService) ToggleTemplateStatus(id, status string) error {
	return s.templateRepo.ToggleTemplateStatus(id, status)
}

// 发起更新
func (s *TemplateService) InitiateUpdate(existingTemplateID, submitterID string) (*models.TemplateRequest, error) {
	existing, err := s.templateRepo.GetTemplateByID(existingTemplateID)
	if err != nil {
		return nil, errors.New("模板不存在")
	}

	req := &models.TemplateRequest{
		ID:                 uuid.New().String(),
		RequestNo:          utils.GenerateTemplateRequestNo(),
		RequestType:        models.TemplateRequestTypeUpdate,
		ExistingTemplateID: &existingTemplateID,
		ContractType:       existing.ContractType,
		Title:              "更新模板：" + existing.Name,
		Status:             models.TemplateRequestStatusPendingApproval,
		CurrentStep:        1,
		SubmitterID:        submitterID,
		SubmittedAt:        time.Now(),
	}

	if err := s.templateRepo.CreateRequest(req); err != nil {
		return nil, err
	}

	return req, nil
}

// 模板统计
func (s *TemplateService) GetTemplateStats() (map[string]interface{}, error) {
	return s.templateRepo.GetTemplateStats()
}
