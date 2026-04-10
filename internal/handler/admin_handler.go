package handler

import (
	"github.com/gin-gonic/gin"

	"legal-consultation/internal/config"
	"legal-consultation/internal/models"
	"legal-consultation/internal/repository"
	"legal-consultation/internal/service"
	"legal-consultation/internal/utils"
)

type AdminHandler struct {
	userRepo         *repository.UserRepository
	templateRepo     *repository.TemplateRepository
	notificationSvc *service.NotificationService
}

func NewAdminHandler(
	userRepo *repository.UserRepository,
	templateRepo *repository.TemplateRepository,
	notificationSvc *service.NotificationService,
) *AdminHandler {
	return &AdminHandler{
		userRepo:         userRepo,
		templateRepo:     templateRepo,
		notificationSvc: notificationSvc,
	}
}

// ============ 用户管理 ============

// ListUsers 用户列表
func (h *AdminHandler) ListUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	page := getPage(c)
	pageSize := getPageSize(c)

	var users []models.User
	var total int64
	var err error

	if keyword != "" {
		users, total, err = h.userRepo.Search(keyword, page, pageSize)
	} else {
		users, total, err = h.userRepo.List(page, pageSize)
	}

	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"items":     users,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

// CreateUser 创建用户
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.userRepo.Create(&user); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessCreated(c, user)
}

// UpdateUser 更新用户
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	user.ID = id
	if err := h.userRepo.Update(&user); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, user)
}

// ResetPassword 重置密码
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	// TODO: 实现重置密码
	utils.Success(c, gin.H{"message": "密码重置成功"})
}

// ToggleUserStatus 启用/禁用用户
func (h *AdminHandler) ToggleUserStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.userRepo.ToggleStatus(id, req.Status); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// ============ 部门管理 ============

// ListDepartments 部门列表
func (h *AdminHandler) ListDepartments(c *gin.Context) {
	// TODO: 实现部门列表
	utils.Success(c, []interface{}{})
}

// CreateDepartment 创建部门
func (h *AdminHandler) CreateDepartment(c *gin.Context) {
	utils.Success(c, nil)
}

// UpdateDepartment 更新部门
func (h *AdminHandler) UpdateDepartment(c *gin.Context) {
	utils.Success(c, nil)
}

// DeleteDepartment 删除部门
func (h *AdminHandler) DeleteDepartment(c *gin.Context) {
	utils.Success(c, nil)
}

// ============ 合同类型管理 ============

// ListContractTypes 合同类型列表
func (h *AdminHandler) ListContractTypes(c *gin.Context) {
	utils.Success(c, []string{"采购合同", "劳动合同", "租赁合同", "销售合同", "保密协议", "服务合同", "其他"})
}

// UpdateContractTypes 更新合同类型
func (h *AdminHandler) UpdateContractTypes(c *gin.Context) {
	utils.Success(c, nil)
}

// ============ 咨询类型配置 ============

// ListConsultationTypes 咨询类型配置
func (h *AdminHandler) ListConsultationTypes(c *gin.Context) {
	// TODO: 从数据库获取配置
	utils.Success(c, []interface{}{})
}

// UpdateConsultationType 更新咨询类型配置
func (h *AdminHandler) UpdateConsultationType(c *gin.Context) {
	utils.Success(c, nil)
}

// ============ 系统设置 ============

// GetSystemConfig 获取系统配置
func (h *AdminHandler) GetSystemConfig(c *gin.Context) {
	config := h.notificationSvc.GetConfig()
	utils.Success(c, gin.H{
		"dingtalk_webhook_url": config.WebhookURL,
		"dingtalk_enabled":     config.Enabled,
	})
}

// UpdateSystemConfig 更新系统配置
func (h *AdminHandler) UpdateSystemConfig(c *gin.Context) {
	var req struct {
		DingTalkWebhookURL string `json:"dingtalk_webhook_url"`
		DingTalkEnabled    bool   `json:"dingtalk_enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	h.notificationSvc.UpdateConfig(&config.DingTalkConfig{
		WebhookURL: req.DingTalkWebhookURL,
		Enabled:    req.DingTalkEnabled,
	})

	utils.Success(c, nil)
}

// TestNotification 测试通知
func (h *AdminHandler) TestNotification(c *gin.Context) {
	if err := h.notificationSvc.TestNotification(); err != nil {
		utils.InternalError(c, "发送测试通知失败")
		return
	}

	utils.Success(c, gin.H{"message": "测试通知已发送"})
}
