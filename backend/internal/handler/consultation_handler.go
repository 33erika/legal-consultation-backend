package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"legal-consultation/internal/middleware"
	"legal-consultation/internal/service"
	"legal-consultation/internal/utils"
)

type ConsultationHandler struct {
	consultationSvc *service.ConsultationService
}

func NewConsultationHandler(consultationSvc *service.ConsultationService) *ConsultationHandler {
	return &ConsultationHandler{consultationSvc: consultationSvc}
}

// Create 创建咨询
func (h *ConsultationHandler) Create(c *gin.Context) {
	var req service.CreateConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user := middleware.GetCurrentUser(c)
	consultation, err := h.consultationSvc.CreateConsultation(&req, user.ID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessCreated(c, gin.H{
		"id":        consultation.ID,
		"ticket_no": consultation.TicketNo,
	})
}

// List 获取咨询列表
func (h *ConsultationHandler) List(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	query := &service.ListConsultationsQuery{
		Role:     user.Role,
		UserID:   user.ID,
		Status:   c.Query("status"),
		Urgency:  c.Query("urgency"),
		Page:     getPage(c),
		PageSize: getPageSize(c),
	}

	consultations, total, err := h.consultationSvc.ListConsultations(query)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":    consultations,
		"total":    total,
		"page":     query.Page,
		"page_size": query.PageSize,
	})
}

// Get 获取咨询详情
func (h *ConsultationHandler) Get(c *gin.Context) {
	id := c.Param("id")

	consultation, err := h.consultationSvc.GetConsultation(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, consultation)
}

// Accept 接单
func (h *ConsultationHandler) Accept(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req service.AcceptConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.consultationSvc.AcceptConsultation(id, user.ID, &req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Reply 回复
func (h *ConsultationHandler) Reply(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req service.ReplyConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.consultationSvc.ReplyConsultation(id, &req, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessCreated(c, nil)
}

// RequestSupplement 要求补充资料
func (h *ConsultationHandler) RequestSupplement(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req service.RequestSupplementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.consultationSvc.RequestSupplement(id, &req, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Close 结案
func (h *ConsultationHandler) Close(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req service.CloseConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.consultationSvc.CloseConsultation(id, &req, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Transfer 变更处理人
func (h *ConsultationHandler) Transfer(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req struct {
		NewHandlerID string `json:"new_handler_id" binding:"required"`
		Reason       string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.consultationSvc.TransferConsultation(id, req.NewHandlerID, req.Reason, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Rate 评价
func (h *ConsultationHandler) Rate(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Rating int `json:"rating" binding:"required,min=1,max=5"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.consultationSvc.RateConsultation(id, req.Rating); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Similar 相似问题
func (h *ConsultationHandler) Similar(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		utils.BadRequest(c, "缺少title参数")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	consultations, err := h.consultationSvc.FindSimilar(title, limit)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, consultations)
}

// GetStats 获取统计
func (h *ConsultationHandler) GetStats(c *gin.Context) {
	stats, err := h.consultationSvc.GetStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

// Search 搜索
func (h *ConsultationHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	status := c.Query("status")
	page := getPage(c)
	pageSize := getPageSize(c)

	consultations, total, err := h.consultationSvc.Search(keyword, category, status, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":    consultations,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func getPage(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

func getPageSize(c *gin.Context) int {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return pageSize
}
