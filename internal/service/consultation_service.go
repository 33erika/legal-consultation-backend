package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"legal-consultation/internal/models"
	"legal-consultation/internal/repository"
	"legal-consultation/internal/utils"
)

type ConsultationService struct {
	consultationRepo *repository.ConsultationRepository
	attachmentRepo   *repository.AttachmentRepository
	userRepo         *repository.UserRepository
	notificationSvc  *NotificationService
}

func NewConsultationService(
	consultationRepo *repository.ConsultationRepository,
	attachmentRepo *repository.AttachmentRepository,
	userRepo *repository.UserRepository,
	notificationSvc *NotificationService,
) *ConsultationService {
	return &ConsultationService{
		consultationRepo: consultationRepo,
		attachmentRepo:   attachmentRepo,
		userRepo:         userRepo,
		notificationSvc:  notificationSvc,
	}
}

// 创建咨询
func (s *ConsultationService) CreateConsultation(req *CreateConsultationRequest, submitterID string) (*models.Consultation, error) {
	consultation := &models.Consultation{
		ID:               uuid.New().String(),
		TicketNo:         utils.GenerateConsultationTicketNo(),
		Title:            req.Title,
		Description:      req.Description,
		Urgency:          req.Urgency,
		Status:           models.ConsultationStatusPending,
		ConsultationType: req.ConsultationType,
		ExtensionData:    req.ExtensionData,
		SubmitterID:      submitterID,
		SubmittedAt:      time.Now(),
	}

	if err := s.consultationRepo.Create(consultation); err != nil {
		return nil, err
	}

	// 发送通知
	s.notificationSvc.NotifyNewConsultation(consultation)

	return consultation, nil
}

type CreateConsultationRequest struct {
	Title            string                   `json:"title" binding:"required,max=100"`
	Description      string                   `json:"description" binding:"required,max=5000"`
	Urgency          string                   `json:"urgency" binding:"required,oneof=normal urgent very_urgent"`
	ConsultationType string                   `json:"consultation_type"`
	ExtensionData    models.JSONType          `json:"extension_data"`
}

// 获取咨询详情
func (s *ConsultationService) GetConsultation(id string) (*models.Consultation, error) {
	consultation, err := s.consultationRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("咨询不存在")
	}

	// 加载回复
	if err := s.loadReplies(consultation); err != nil {
		return nil, err
	}

	// 加载附件
	if err := s.loadAttachments(consultation); err != nil {
		return nil, err
	}

	return consultation, nil
}

func (s *ConsultationService) loadReplies(c *models.Consultation) error {
	var replies []models.ConsultationReply
	// 这里简化处理，实际应该从数据库加载
	c.Replies = replies
	return nil
}

func (s *ConsultationService) loadAttachments(c *models.Consultation) error {
	attachments, err := s.attachmentRepo.ListConsultationAttachments(c.ID)
	if err != nil {
		return err
	}

	var result []models.ConsultationAttachment
	for _, a := range attachments {
		result = append(result, a)
	}
	c.Attachments = result
	return nil
}

// 接单
func (s *ConsultationService) AcceptConsultation(id, handlerID string, req *AcceptConsultationRequest) error {
	consultation, err := s.consultationRepo.GetByID(id)
	if err != nil {
		return errors.New("咨询不存在")
	}

	if consultation.Status != models.ConsultationStatusPending {
		return errors.New("该咨询已被接单")
	}

	now := time.Now()
	consultation.HandlerID = &handlerID
	consultation.Status = models.ConsultationStatusProcessing
	consultation.InternalCategory = req.InternalCategory
	consultation.ComplexSubCategory = req.ComplexSubCategory
	consultation.AcceptedAt = &now

	if err := s.consultationRepo.Update(consultation); err != nil {
		return err
	}

	// 发送通知
	s.notificationSvc.NotifyConsultationAccepted(consultation)

	return nil
}

type AcceptConsultationRequest struct {
	InternalCategory   string `json:"internal_category" binding:"required,oneof=simple complex"`
	ComplexSubCategory string `json:"complex_sub_category"`
}

// 回复咨询
func (s *ConsultationService) ReplyConsultation(id string, req *ReplyConsultationRequest, authorID string) error {
	consultation, err := s.consultationRepo.GetByID(id)
	if err != nil {
		return errors.New("咨询不存在")
	}

	if consultation.Status == models.ConsultationStatusClosed {
		return errors.New("该咨询已结案，无法回复")
	}

	reply := &models.ConsultationReply{
		ID:             uuid.New().String(),
		ConsultationID: id,
		ReplyType:      models.ReplyTypeLegal,
		Content:        req.Content,
		AuthorID:       authorID,
		CreatedAt:      time.Now(),
	}

	// 更新咨询状态
	consultation.Status = models.ConsultationStatusReplied
	if consultation.FirstRepliedAt == nil {
		consultation.FirstRepliedAt = &reply.CreatedAt
	}

	if err := s.consultationRepo.Update(consultation); err != nil {
		return err
	}

	// 发送通知
	s.notificationSvc.NotifyConsultationReplied(consultation)

	return nil
}

type ReplyConsultationRequest struct {
	Content string `json:"content" binding:"required"`
}

// 要求补充资料
func (s *ConsultationService) RequestSupplement(id string, req *RequestSupplementRequest, authorID string) error {
	consultation, err := s.consultationRepo.GetByID(id)
	if err != nil {
		return errors.New("咨询不存在")
	}

	// 状态保持不变（不退回待处理）
	if err := s.consultationRepo.Update(consultation); err != nil {
		return err
	}

	// 发送通知
	s.notificationSvc.NotifyConsultationReplied(consultation)

	return nil
}

type RequestSupplementRequest struct {
	Message string `json:"message" binding:"required"`
}

// 结案
func (s *ConsultationService) CloseConsultation(id string, req *CloseConsultationRequest, operatorID string) error {
	consultation, err := s.consultationRepo.GetByID(id)
	if err != nil {
		return errors.New("咨询不存在")
	}

	if consultation.Status == models.ConsultationStatusClosed {
		return errors.New("该咨询已结案")
	}

	now := time.Now()
	consultation.Status = models.ConsultationStatusClosed
	consultation.ClosedAt = &now

	// 如果提供了新的分类，使用新分类
	if req.InternalCategory != "" {
		consultation.InternalCategory = req.InternalCategory
	}
	if req.ComplexSubCategory != "" {
		consultation.ComplexSubCategory = req.ComplexSubCategory
	}

	if err := s.consultationRepo.Update(consultation); err != nil {
		return err
	}

	// 发送通知
	s.notificationSvc.NotifyConsultationClosed(consultation)

	return nil
}

type CloseConsultationRequest struct {
	InternalCategory   string `json:"internal_category"`
	ComplexSubCategory string `json:"complex_sub_category"`
}

// 变更处理人
func (s *ConsultationService) TransferConsultation(id, newHandlerID, reason, operatorID string) error {
	consultation, err := s.consultationRepo.GetByID(id)
	if err != nil {
		return errors.New("咨询不存在")
	}

	if consultation.HandlerID != nil && *consultation.HandlerID == newHandlerID {
		return errors.New("无法转让给同一人")
	}

	oldHandlerID := consultation.HandlerID
	consultation.HandlerID = &newHandlerID

	if err := s.consultationRepo.Update(consultation); err != nil {
		return err
	}

	// 发送通知
	s.notificationSvc.NotifyConsultationTransferred(consultation, oldHandlerID)

	return nil
}

// 评价
func (s *ConsultationService) RateConsultation(id string, rating int) error {
	if rating < 1 || rating > 5 {
		return errors.New("评分必须在1-5之间")
	}

	consultation, err := s.consultationRepo.GetByID(id)
	if err != nil {
		return errors.New("咨询不存在")
	}

	consultation.Rating = &rating
	return s.consultationRepo.Update(consultation)
}

// 列表
func (s *ConsultationService) ListConsultations(query *ListConsultationsQuery) ([]models.Consultation, int64, error) {
	switch query.Role {
	case "staff":
		return s.consultationRepo.ListBySubmitter(query.UserID, query.Status, query.Page, query.PageSize)
	case "legal_staff":
		if query.Status == models.ConsultationStatusPending {
			return s.consultationRepo.ListPending(query.Page, query.PageSize, query.Urgency)
		}
		return s.consultationRepo.ListByHandler(query.UserID, query.Status, query.Page, query.PageSize)
	default:
		return nil, 0, errors.New("无效的角色")
	}
}

type ListConsultationsQuery struct {
	Role     string `json:"role"`
	UserID   string `json:"user_id"`
	Status   string `json:"status"`
	Urgency  string `json:"urgency"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// 统计
func (s *ConsultationService) GetStats() (map[string]interface{}, error) {
	return s.consultationRepo.GetStats()
}

// 搜索
func (s *ConsultationService) Search(keyword, category, status string, page, pageSize int) ([]models.Consultation, int64, error) {
	return s.consultationRepo.Search(keyword, category, status, page, pageSize)
}

// 相似问题
func (s *ConsultationService) FindSimilar(title string, limit int) ([]models.Consultation, error) {
	return s.consultationRepo.FindSimilar(title, limit)
}
