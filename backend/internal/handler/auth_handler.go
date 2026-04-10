package handler

import (
	"github.com/gin-gonic/gin"

	"legal-consultation/internal/models"
	"legal-consultation/internal/service"
	"legal-consultation/internal/utils"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	resp, err := h.authSvc.Login(&req)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	utils.Success(c, nil)
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.Unauthorized(c, "请先登录")
		return
	}
	utils.Success(c, user)
}

// CreateUser 创建用户 (管理员)
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	// 密码加密
	// TODO: 这里应该调用 password service

	if err := h.authSvc.CreateUser(&user); err != nil {
		utils.InternalError(c, "创建用户失败")
		return
	}

	utils.SuccessCreated(c, user)
}
