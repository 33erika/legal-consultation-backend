package handler

import (
	"github.com/gin-gonic/gin"

	"legal-consultation/internal/middleware"
	"legal-consultation/internal/service"
	"legal-consultation/internal/utils"
)

type TemplateHandler struct {
	templateSvc *service.TemplateService
}

func NewTemplateHandler(templateSvc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{templateSvc: templateSvc}
}

// ============ 模板申请 ============

// CreateRequest 创建申请
func (h *TemplateHandler) CreateRequest(c *gin.Context) {
	var req service.CreateTemplateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	user := middleware.GetCurrentUser(c)
	templateReq, err := h.templateSvc.CreateTemplateRequest(&req, user.ID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessCreated(c, gin.H{
		"id":         templateReq.ID,
		"request_no": templateReq.RequestNo,
	})
}

// ListMyRequests 我的申请列表
func (h *TemplateHandler) ListMyRequests(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	requests, total, err := h.templateSvc.ListMyRequests(user.ID, getPage(c), getPageSize(c))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     requests,
		"total":    total,
		"page":     getPage(c),
		"page_size": getPageSize(c),
	})
}

// GetRequest 获取申请详情
func (h *TemplateHandler) GetRequest(c *gin.Context) {
	id := c.Param("id")

	req, err := h.templateSvc.GetTemplateRequest(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, req)
}

// Approve L1审批
func (h *TemplateHandler) Approve(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req service.ApproveTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.templateSvc.ApproveRequest(id, &req, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Draft 拟写模板
func (h *TemplateHandler) Draft(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req service.DraftTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.templateSvc.DraftTemplate(id, &req, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// SaveDraft 保存草稿
func (h *TemplateHandler) SaveDraft(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	if err := h.templateSvc.SaveDraft(id, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Review 审核
func (h *TemplateHandler) Review(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetCurrentUser(c)

	var req service.ReviewTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.templateSvc.ReviewTemplate(id, &req, user.ID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// ListPendingApproval 待审批列表
func (h *TemplateHandler) ListPendingApproval(c *gin.Context) {
	requests, total, err := h.templateSvc.ListPendingApproval(getPage(c), getPageSize(c))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     requests,
		"total":    total,
		"page":     getPage(c),
		"page_size": getPageSize(c),
	})
}

// ListPendingDraft 待拟写列表
func (h *TemplateHandler) ListPendingDraft(c *gin.Context) {
	requests, total, err := h.templateSvc.ListPendingDraft(getPage(c), getPageSize(c))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     requests,
		"total":    total,
		"page":     getPage(c),
		"page_size": getPageSize(c),
	})
}

// ListPendingReview 待审核列表
func (h *TemplateHandler) ListPendingReview(c *gin.Context) {
	requests, total, err := h.templateSvc.ListPendingReview(getPage(c), getPageSize(c))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     requests,
		"total":    total,
		"page":     getPage(c),
		"page_size": getPageSize(c),
	})
}

// GetRequestStats 获取申请统计
func (h *TemplateHandler) GetRequestStats(c *gin.Context) {
	stats, err := h.templateSvc.GetRequestStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

// ============ 合同模板库 ============

// ListTemplates 模板列表
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	templates, total, err := h.templateSvc.ListTemplates(
		c.Query("contract_type"),
		c.Query("keyword"),
		getPage(c),
		getPageSize(c),
	)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     templates,
		"total":    total,
		"page":     getPage(c),
		"page_size": getPageSize(c),
	})
}

// GetTemplate 获取模板详情
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")

	template, err := h.templateSvc.GetTemplate(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, template)
}

// DownloadTemplate 下载模板
func (h *TemplateHandler) DownloadTemplate(c *gin.Context) {
	id := c.Param("id")

	if err := h.templateSvc.DownloadTemplate(id); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetTemplateVersions 模板版本历史
func (h *TemplateHandler) GetTemplateVersions(c *gin.Context) {
	id := c.Param("id")

	versions, err := h.templateSvc.GetTemplateVersions(id)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, versions)
}

// CompareVersions 版本对比
func (h *TemplateHandler) CompareVersions(c *gin.Context) {
	versionAID := c.Query("version_a_id")
	versionBID := c.Query("version_b_id")

	if versionAID == "" || versionBID == "" {
		utils.BadRequest(c, "缺少版本ID参数")
		return
	}

	versionA, versionB, err := h.templateSvc.CompareVersions(versionAID, versionBID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"version_a": versionA,
		"version_b": versionB,
	})
}

// ToggleTemplateStatus 禁用/启用模板
func (h *TemplateHandler) ToggleTemplateStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required,oneof=published disabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.templateSvc.ToggleTemplateStatus(id, req.Status); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// InitiateUpdate 发起更新
func (h *TemplateHandler) InitiateUpdate(c *gin.Context) {
	var req struct {
		TemplateID string `json:"template_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	user := middleware.GetCurrentUser(c)
	templateReq, err := h.templateSvc.InitiateUpdate(req.TemplateID, user.ID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessCreated(c, gin.H{
		"id": templateReq.ID,
	})
}

// GetTemplateStats 获取模板统计
func (h *TemplateHandler) GetTemplateStats(c *gin.Context) {
	stats, err := h.templateSvc.GetTemplateStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}
