package handler

import (
	"github.com/gin-gonic/gin"

	"legal-consultation/internal/middleware"
	"legal-consultation/internal/service"
	"legal-consultation/internal/utils"
)

type LegalHandler struct {
	consultationSvc *service.ConsultationService
	templateSvc     *service.TemplateService
}

func NewLegalHandler(
	consultationSvc *service.ConsultationService,
	templateSvc *service.TemplateService,
) *LegalHandler {
	return &LegalHandler{
		consultationSvc: consultationSvc,
		templateSvc:     templateSvc,
	}
}

// Dashboard 工作台统计
func (h *LegalHandler) Dashboard(c *gin.Context) {
	consultationStats, err := h.consultationSvc.GetStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	templateStats, err := h.templateSvc.GetRequestStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"consultation_stats": consultationStats,
		"template_stats":    templateStats,
	})
}

// ConsultationPool 咨询池
func (h *LegalHandler) ConsultationPool(c *gin.Context) {
	urgency := c.Query("urgency")
	page := getPage(c)
	pageSize := getPageSize(c)

	consultations, total, err := h.consultationSvc.ListConsultations(&service.ListConsultationsQuery{
		Role:     "legal_staff",
		Status:   "pending",
		Urgency:  urgency,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     consultations,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

// MyTasks 我的待办
func (h *LegalHandler) MyTasks(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	status := c.Query("status")
	page := getPage(c)
	pageSize := getPageSize(c)

	consultations, total, err := h.consultationSvc.ListConsultations(&service.ListConsultationsQuery{
		Role:     "legal_staff",
		UserID:   user.ID,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     consultations,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

// StaffList 法务专员列表
func (h *LegalHandler) StaffList(c *gin.Context) {
	// TODO: 从 user repo 获取法务专员列表
	utils.Success(c, []interface{}{})
}
